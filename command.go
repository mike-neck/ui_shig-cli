package main

import "github.com/urfave/cli/v2"

type Command interface {
	Execute(config UiShigConfig, voices []Voice) error
}

type FilePath string

type OptionType interface {
	int |
		string |
		FilePath |
		bool
}

type Ref[T OptionType] *T

type Option[T OptionType] struct {
	Name        string
	Aliases     []string
	Description string
	Required    bool
	Reference   Ref[T]
}

type IntOption Option[int]

type StringOption Option[string]

type FileOption Option[FilePath]

type BoolOption Option[bool]

type UserOrder struct {
	Name                string
	Description         string
	ArgumentDescription *string
	IntOptions          []IntOption
	StringOptions       []StringOption
	FileOptions         []FileOption
	BoolOptions         []BoolOption
	ConstructCommand    func(uo UserOrder, args []string) (Command, error)
}

type Reflector[C any, T OptionType] struct {
	CLIType *C
	Option  *T
}

func (order UserOrder) toCLICommand(config UiShigConfig, voices []Voice) *cli.Command {
	flags := make([]cli.Flag, 0)
	for _, option := range order.IntOptions {
		flags = append(flags, &cli.IntFlag{
			Name:        option.Name,
			Aliases:     option.Aliases,
			Usage:       option.Description,
			Required:    option.Required,
			Destination: option.Reference,
		})
	}
	for _, option := range order.StringOptions {
		flags = append(flags, &cli.StringFlag{
			Name:        option.Name,
			Aliases:     option.Aliases,
			Usage:       option.Description,
			Required:    option.Required,
			Destination: option.Reference,
		})
	}
	reflectors := make([]Reflector[cli.Path, FilePath], 0)
	for _, option := range order.FileOptions {
		var cliPath cli.Path
		fileOption := option
		reflector := Reflector[cli.Path, FilePath]{
			CLIType: &cliPath,
			Option:  fileOption.Reference,
		}
		reflectors = append(reflectors, reflector)
		flags = append(flags, &cli.PathFlag{
			Name:        option.Name,
			Aliases:     option.Aliases,
			Usage:       option.Description,
			Required:    option.Required,
			Destination: reflector.CLIType,
		})
	}
	for _, option := range order.BoolOptions {
		flags = append(flags, &cli.BoolFlag{
			Name:        option.Name,
			Aliases:     option.Aliases,
			Usage:       option.Description,
			Required:    option.Required,
			Destination: option.Reference,
		})
	}
	var argsUsage string
	if order.ArgumentDescription != nil {
		argsUsage = *order.ArgumentDescription
	}
	return &cli.Command{
		Name:        order.Name,
		Usage:       order.Description,
		UsageText:   order.Description,
		Description: order.Description,
		Args:        order.ArgumentDescription != nil,
		ArgsUsage:   argsUsage,
		Flags:       flags,
		Action: func(context *cli.Context) error {
			for _, r := range reflectors {
				filePath := FilePath(*r.CLIType)
				if filePath != "" {
					r.Option = &filePath
				}
			}
			command, err := order.ConstructCommand(order, context.Args().Slice())
			if err != nil {
				return err
			}
			return command.Execute(config, voices)
		},
	}
}
