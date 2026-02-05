extends Node

# Global Configuration
var claude_api_key: String = ""
var openai_api_key: String = ""

# Mock Switch - Set to true by default for safety/testing
var use_mocks: bool = true

func _ready():
	_load_env()
	print("Global: Ready. API Keys loaded present? ", claude_api_key != "", openai_api_key != "")
	print("Global: Use Mocks? ", use_mocks)

func _load_env():
	# Simple .env loader
	var file = FileAccess.open("res://.env", FileAccess.READ)
	if file:
		while not file.eof_reached():
			var line = file.get_line()
			if line.strip_edges().is_empty() or line.begins_with("#"):
				continue
			var parts = line.split("=", true, 1)
			if parts.size() == 2:
				var key = parts[0].strip_edges()
				var value = parts[1].strip_edges()
				if key == "CLAUDE_API_KEY":
					claude_api_key = value
				elif key == "OPENAI_API_KEY":
					openai_api_key = value
	else:
		print("Global: No .env file found. Ensure keys are set or use mocks.")
