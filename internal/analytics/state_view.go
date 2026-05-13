package analytics

import (
	"fmt"
	"math/big"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/vivlilv/go_nft_trader/internal/domain"
)

var (
	titleStyle  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("214"))
	headerStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("39")).Align(lipgloss.Center)
	cellStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("252")).Padding(0, 1)
	dimStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Padding(0, 1)
	greenStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("82")).Padding(0, 1)
	borderStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
)

type tickMsg time.Time

type model struct {
	stateManager *domain.StateManager
	snapshot     domain.State
}

func NewModel(stateManager *domain.StateManager) model {
	state := stateManager.GetState()
	return model{stateManager: stateManager, snapshot: *state}
}

func tickCmd() tea.Cmd {
	return tea.Tick(200*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m model) Init() tea.Cmd {
	return tickCmd()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tickMsg:
		m.snapshot = m.stateManager.Snapshot()
		return m, tickCmd()

	case tea.KeyMsg:
		if msg.String() == "q" {
			return m, tea.Quit
		}
	}
	return m, nil
}

var weiPerEth = new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil))

func weiToEth(wei *big.Int) string {
	if wei == nil {
		return "—"
	}
	eth := new(big.Float).Quo(new(big.Float).SetInt(wei), weiPerEth)
	return fmt.Sprintf("%.4f ETH", eth)
}

func (m model) View() string {
	rows := [][]string{}
	for id, item := range m.snapshot.Items {
		myEth := weiToEth(item.MyOffer.PriceWei)
		myCell := cellStyle.Render(myEth)
		if item.MyOffer.PriceWei != nil {
			myCell = greenStyle.Render(myEth)
		}
		rows = append(rows, []string{
			string(id),
			weiToEth(item.TopOffer.PriceWei),
			myCell,
		})
	}

	t := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(borderStyle).
		Headers("NftID", "TopOffer", "MyOffer").
		StyleFunc(func(row, col int) lipgloss.Style {
			if row == table.HeaderRow {
				return headerStyle
			}
			return cellStyle
		}).
		Rows(rows...)

	return titleStyle.Render("NFT Trader — Live Positions") + "\n\n" +
		t.Render() + "\n\n" +
		dimStyle.Render(" Press q to quit.")
}

func Run(sm *domain.StateManager) error {
	p := tea.NewProgram(NewModel(sm), tea.WithAltScreen())
	_, err := p.Run()
	return err
}
