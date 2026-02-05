extends Control

@onready var chat_log = $UI/ChatLog
@onready var input_field = $UI/InputRow/InputElement
@onready var send_button = $UI/InputRow/SendButton
@onready var ai_service = $AIService
@onready var tts_service = $TTSService
@onready var character = $CharacterContainer/SubViewport/SpineCharacter
@onready var audio_player = $AudioStreamPlayer

func _ready():
	send_button.pressed.connect(_on_send_pressed)
	input_field.text_submitted.connect(_on_send_pressed.bind(""))
	
	ai_service.response_received.connect(_on_ai_response)
	ai_service.error_occurred.connect(_on_error)
	
	tts_service.audio_received.connect(_on_audio_received)
	tts_service.error_occurred.connect(_on_error)
	
	audio_player.finished.connect(_on_audio_finished)
	
	_log_message("System", "Ready. Mock Mode: " + str(Global.use_mocks))
	character.play_idle()

func _on_send_pressed(_text = ""):
	var text = input_field.text
	if text.strip_edges().is_empty():
		return
		
	_log_message("User", text)
	input_field.text = ""
	input_field.editable = false
	
	character.play_listen()
	ai_service.send_message(text)

func _on_ai_response(text: String):
	_log_message("AI", text)
	tts_service.speak(text)

func _on_audio_received(stream: AudioStreamMP3):
	if stream:
		audio_player.stream = stream
		audio_player.play()
		character.play_speak()
	else:
		# Mock audio
		print("Main: Mock audio received (silence/null). Simulating playback time.")
		character.play_speak()
		# Fake duration for mock
		get_tree().create_timer(3.0).timeout.connect(_on_audio_finished)

func _on_audio_finished():
	character.play_idle()
	input_field.editable = true
	input_field.grab_focus()

func _on_error(message: String):
	_log_message("Error", message)
	input_field.editable = true
	character.play_idle()

func _log_message(sender: String, message: String):
	var color = "#ffffff"
	if sender == "User": color = "#aaffaa"
	elif sender == "AI": color = "#aaaaff"
	elif sender == "Error": color = "#ffaaaa"
	
	chat_log.append_text("[color=%s][b]%s:[/b] %s[/color]\n" % [color, sender, message])
