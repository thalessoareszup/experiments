.PHONY: install help

help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  install    Install the skills command to ~/.local/bin/skills"
	@echo "  help       Show this help message"

install:
	@mkdir -p "$(HOME)/.local/bin"
	@ln -sf "$(CURDIR)/skills" "$(HOME)/.local/bin/skills"
	@chmod +x "$(CURDIR)/skills"
	@echo "✓ Installed skills command to ~/.local/bin/skills"
	@echo ""
	@echo "Make sure ~/.local/bin is in your PATH:"
	@echo '  export PATH="$$HOME/.local/bin:$$PATH"'
