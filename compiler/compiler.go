package compiler

import (
	"context"
	"io"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/sirupsen/logrus"
)

const (
	FirmwareName = "ex_dactyl_v2_custom.hex"

	cmdTimeout  = time.Second * 2
	targetImage = "lschyi/qmk_compiler.slim"
	//targetImage  = "lschyi/qmk_compiler:20200110"
	keymapPath   = "/root/qmk_firmware/keyboards/ex_dactyl_v2/keymaps/custom"
	keymapName   = "keymap.c"
	firmwarePath = "/root/qmk_firmware/" + FirmwareName
)

type Compiler struct {
	client         *client.Client
	containerID    string
	cancelCopyBack context.CancelFunc
	firmwareReader io.ReadCloser
}

func NewCompiler() (*Compiler, error) {
	c, err := client.NewEnvClient()
	if err != nil {
		return nil, err
	}

	res := &Compiler{client: c}
	if err := res.setupContainer(); err != nil {
		return nil, err
	}
	return res, nil
}

func (c *Compiler) Compile(keymap io.Reader) (io.Reader, error) {
	if err := c.putKeymap(keymap); err != nil {
		return nil, err
	}
	if err := c.startContainer(); err != nil {
		return nil, err
	}
	if err := c.waitContainer(); err != nil {
		return nil, err
	}
	return c.getFirmware()
}

func (c *Compiler) Close() {
	if c.containerID != "" {
		c.removeContainer()
	}
	if c.firmwareReader != nil {
		c.firmwareReader.Close()
	}
	if c.cancelCopyBack != nil {
		c.cancelCopyBack()
	}
	c.client.Close()
}

func (c *Compiler) setupContainer() error {
	config := &container.Config{Image: targetImage}
	ctx, cancel := context.WithTimeout(context.Background(), cmdTimeout)
	defer cancel()
	createBody, err := c.client.ContainerCreate(ctx, config, nil, nil, "")
	if err != nil {
		return err
	}
	c.containerID = createBody.ID
	return nil
}

func (c *Compiler) putKeymap(source io.Reader) error {
	ctx, cancel := context.WithTimeout(context.Background(), cmdTimeout)
	defer cancel()
	tarReader, err := genTarKeymap(source)
	if err != nil {
		return err
	}
	config := types.CopyToContainerOptions{}
	return c.client.CopyToContainer(ctx, c.containerID, keymapPath, tarReader, config)
}

func (c *Compiler) startContainer() error {
	ctx, cancel := context.WithTimeout(context.Background(), cmdTimeout)
	defer cancel()
	return c.client.ContainerStart(ctx, c.containerID, types.ContainerStartOptions{})
}

func (c *Compiler) waitContainer() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if _, err := c.client.ContainerWait(ctx, c.containerID); err != nil {
		return err
	}
	return nil
}

func (c *Compiler) getFirmware() (io.Reader, error) {
	ctx, cancel := context.WithTimeout(context.Background(), cmdTimeout)
	c.cancelCopyBack = cancel
	reader, _, err := c.client.CopyFromContainer(ctx, c.containerID, firmwarePath)
	if err != nil {
		return nil, err
	}
	c.firmwareReader = reader
	return unTarFirmware(reader)
}

func (c *Compiler) removeContainer() {
	ctx, cancel := context.WithTimeout(context.Background(), cmdTimeout)
	defer cancel()
	config := types.ContainerRemoveOptions{Force: true}
	if err := c.client.ContainerRemove(ctx, c.containerID, config); err != nil {
		logrus.WithError(err).Error("encounter error while removing container")
	}
}
