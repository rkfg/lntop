package ui

import (
	"context"

	"github.com/jroimartin/gocui"
	"github.com/pkg/errors"

	"github.com/edouardparis/lntop/config"
	"github.com/edouardparis/lntop/events"
	"github.com/edouardparis/lntop/logging"
	"github.com/edouardparis/lntop/network"
)

func Run(ctx context.Context, cfg config.Config, logger logging.Logger, net network.Network, sub chan *events.Event) error {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		return err
	}
	defer g.Close()

	g.Cursor = true
	ctrl := newController(cfg, logger, net)
	err = ctrl.SetModels(ctx)
	if err != nil {
		return err
	}

	g.SetManagerFunc(ctrl.layout)

	err = setKeyBinding(ctrl, g)
	if err != nil {
		return err
	}

	go ctrl.Listen(ctx, g, sub)

	err = g.MainLoop()
	close(sub)

	return errors.WithStack(err)
}
