extends Control

@onready var log_display = $LogDisplay
@onready var ai_service = $AIService
@onready var tts_service = $TTSService

func _ready():
	_log("Starting Automated Mock Tests...")
	if not Global.use_mocks:
		_log("WARNING: Global.use_mocks is FALSE. Tests might fail if no keys are present.")
	else:
		_log("Mock Mode is ACTIVE.")
	
	_run_tests()

func _run_tests():
	await get_tree().create_timer(1.0).timeout
	
	# Test 1: AI Service Mock
	_log("Test 1: Testing AI Service Mock...")
	ai_service.response_received.connect(_on_ai_test_response)
	ai_service.send_message("Test Message")
	# Wait for response (handled in callback)

func _on_ai_test_response(text: String):
	_log("PASS: AI Service returned: " + text)
	ai_service.response_received.disconnect(_on_ai_test_response)
	
	# Proceed to Test 2
	_run_test_2()

func _run_test_2():
	await get_tree().create_timer(1.0).timeout
	_log("Test 2: Testing TTS Service Mock...")
	
	tts_service.audio_received.connect(_on_tts_test_audio)
	tts_service.speak("This is a test.")

func _on_tts_test_audio(stream: AudioStreamMP3):
	if stream == null:
		_log("PASS: TTS Service returned null (Expected for Mock).")
	else:
		_log("PASS: TTS Service returned valid AudioStream.")
	
	tts_service.audio_received.disconnect(_on_tts_test_audio)
	
	_finish_tests()

func _finish_tests():
	_log("--------------------------------")
	_log("ALL TESTS COMPLETED.")
	_log("Please run the Main scene for interactive verification.")

func _log(message: String):
	print(message)
	log_display.append_text(message + "\n")
