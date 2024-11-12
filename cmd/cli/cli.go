package cli

import (
	"bufio"
	"fmt"
	"github.com/fcopulgar/stock-manager-go/services"
	"io"
	"strings"
)

type CLI struct {
	portfolioService services.PortfolioServiceInterface
	reader           *bufio.Reader
	writer           io.Writer
}

func NewCLI(portfolioService services.PortfolioServiceInterface, input io.Reader, output io.Writer) *CLI {
	return &CLI{
		portfolioService: portfolioService,
		reader:           bufio.NewReader(input),
		writer:           output,
	}
}

func (cli *CLI) Run() {
	for {
		fmt.Fprintln(cli.writer, "\nSelect an option:")
		fmt.Fprintln(cli.writer, "1. View portfolios")
		fmt.Fprintln(cli.writer, "2. Create portfolio manually")
		fmt.Fprintln(cli.writer, "3. Create random portfolio")
		fmt.Fprintln(cli.writer, "4. Exit")

		fmt.Fprint(cli.writer, "Option: ")
		input, err := cli.reader.ReadString('\n')
		if err != nil {
			fmt.Fprintf(cli.writer, "Error reading input: %v\n", err)
			continue
		}
		input = strings.TrimSpace(input)

		switch input {
		case "1":
			cli.viewPortfolios()
		case "2":
			cli.createPortfolioManual()
		case "3":
			cli.createPortfolioRandom()
		case "4":
			fmt.Fprintln(cli.writer, "Exiting...")
			return
		default:
			fmt.Fprintln(cli.writer, "Invalid option.")
		}
	}
}
