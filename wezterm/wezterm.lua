local wezterm = require("wezterm")
local theme = wezterm.plugin.require("https://github.com/neapsix/wezterm").main
return {
	animation_fps = 60,
	audible_bell = "Disabled",
	automatically_reload_config = true,
	colors = theme.colors(),
	enable_tab_bar = false,
	font = wezterm.font({ family = "MonoLisa Nerd Font", harfbuzz_features = { "ss02", "ss18" } }),
	font_rules = {
		-- normal-intensity-and-italic
		{
			intensity = "Normal",
			italic = true,
			font = wezterm.font_with_fallback({
				family = "MonoLisa Nerd Font",
				-- weight = "DemiLight",
				italic = true,
			}),
		},
	},
	font_size = 18,
	initial_cols = 127,
	initial_rows = 37,
	macos_window_background_blur = 34,
	window_background_opacity = 0.9,
	window_close_confirmation = "NeverPrompt",
	window_decorations = "NONE",
	window_padding = {
		left = 0,
		right = 0,
		top = 0,
		bottom = 0,
	},
}
