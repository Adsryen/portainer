package migrator

import (
	"github.com/rs/zerolog/log"
)

func (m *Migrator) migrateDBVersionToDB33() error {
	log.Info().Msg("updating settings")

	return m.migrateSettingsToDB33()
}

func (m *Migrator) migrateSettingsToDB33() error {
	log.Info().Msg("setting default kubctl shell")
	settings, err := m.settingsService.Settings()
	if err != nil {
		return err
	}

	log.Info().Msg("setting default kubectl shell image")
	settings.KubectlShellImage = *m.flags.KubectlShellImage

	return m.settingsService.UpdateSettings(settings)
}
