package models

import tea "github.com/charmbracelet/bubbletea"

type MainModel struct{}

func NewMainModel() *MainModel

func (m MainModel) Init() tea.Cmd

func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd)

func (m MainModel) View() string
