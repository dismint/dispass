package uconst

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/spf13/viper"
)

func LoadConfig() {
	viper.SetConfigName("dispass")

	viper.AddConfigPath("$HOME/dispass")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		log.Errorf("fatal error reading config file: %v", err)
	}

	// colors
	viper.SetDefault("colors.light.symbol", lostCentury10)
	viper.SetDefault("colors.dark.symbol", lostCentury12)
	viper.SetDefault("colors.light.text", lostCentury10)
	viper.SetDefault("colors.dark.text", lostCentury11)
	viper.SetDefault("colors.light.help_key", lostCentury10)
	viper.SetDefault("colors.dark.help_key", lostCentury15)
	viper.SetDefault("colors.light.help_desc", lostCentury10)
	viper.SetDefault("colors.dark.help_desc", lostCentury16)
	viper.SetDefault("colors.light.help_sep", lostCentury10)
	viper.SetDefault("colors.dark.help_sep", lostCentury14)
	viper.SetDefault("colors.light.border", lostCentury10)
	viper.SetDefault("colors.dark.border", lostCentury12)
	viper.SetDefault("colors.light.message_error", lostCentury10)
	viper.SetDefault("colors.dark.message_error", lostCentury2)
	viper.SetDefault("colors.light.message_success", lostCentury10)
	viper.SetDefault("colors.dark.message_success", lostCentury12)
	viper.SetDefault("colors.light.message_notif", lostCentury10)
	viper.SetDefault("colors.dark.message_notif", lostCentury13)

	// set styles
	SymbolStyle = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{
		Light: viper.GetString("colors.light.symbol"),
		Dark:  viper.GetString("colors.dark.symbol"),
	})
	TextStyle = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{
		Light: viper.GetString("colors.light.text"),
		Dark:  viper.GetString("colors.dark.text"),
	})
	HelpKeyStyle = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{
		Light: viper.GetString("colors.light.help_key"),
		Dark:  viper.GetString("colors.dark.help_key"),
	})
	HelpDescStyle = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{
		Light: viper.GetString("colors.light.help_desc"),
		Dark:  viper.GetString("colors.dark.help_desc"),
	})
	HelpSeparatorStyle = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{
		Light: viper.GetString("colors.light.help_sep"),
		Dark:  viper.GetString("colors.dark.help_sep"),
	})
	BorderColor = lipgloss.AdaptiveColor{
		Light: viper.GetString("colors.light.border"),
		Dark:  viper.GetString("colors.dark.border"),
	}

	// message styles
	MessageBaseStyle := lipgloss.NewStyle().Padding(0, 1)
	MessageLevelErrorStyle = MessageBaseStyle.Foreground(lipgloss.AdaptiveColor{
		Light: viper.GetString("colors.light.message_error"),
		Dark:  viper.GetString("colors.dark.message_error"),
	})
	MessageLevelSuccessStyle = MessageBaseStyle.Foreground(lipgloss.AdaptiveColor{
		Light: viper.GetString("colors.light.message_success"),
		Dark:  viper.GetString("colors.dark.message_success"),
	})
	MessageLevelNotifStyle = MessageBaseStyle.Foreground(lipgloss.AdaptiveColor{
		Light: viper.GetString("colors.light.message_notif"),
		Dark:  viper.GetString("colors.dark.message_notif"),
	})

	HelpStyles = help.Styles{
		ShortKey:       HelpKeyStyle,
		ShortDesc:      HelpDescStyle,
		ShortSeparator: HelpSeparatorStyle,
		FullKey:        HelpKeyStyle,
		FullDesc:       HelpDescStyle,
		FullSeparator:  HelpSeparatorStyle,
	}
}
