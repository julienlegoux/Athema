extends Node
class_name TTSService

signal audio_received(audio_stream: AudioStreamMP3)
signal error_occurred(message: String)

const API_URL = "https://api.openai.com/v1/audio/speech"
const MODEL = "tts-1"
const VOICE = "alloy"

func speak(text: String):
	if Global.use_mocks:
		_mock_speak(text)
	else:
		_real_speak(text)

func _mock_speak(text: String):
	print("TTSService (Mock): Generating audio for -> ", text)
	await get_tree().create_timer(0.5).timeout
	
	# In a real mock scenario without external assets, we trigger the signal 
	# effectively skipping actual audio playback logic or emitting a dummy stream if needed.
	# For visualization, we might just emit a signal that the 'Main' script interprets.
	# But to be robust, let's try to generate a tiny silence MP3 or just emit null and handle it.
	
	# Using a sine wave generator is complex for MP3. 
	# We will emit 'null' and let the player handle "mock playback" visual cues.
	audio_received.emit(null) 

func _real_speak(text: String):
	if Global.openai_api_key.is_empty():
		error_occurred.emit("Missing OpenAI API Key")
		return

	var http_request = HTTPRequest.new()
	add_child(http_request)
	http_request.request_completed.connect(_on_request_completed.bind(http_request))

	var headers = [
		"Authorization: Bearer " + Global.openai_api_key,
		"Content-Type: application/json"
	]
	
	var body = {
		"model": MODEL,
		"input": text,
		"voice": VOICE
	}
	
	var json_body = JSON.stringify(body)
	var error = http_request.request(API_URL, headers, HTTPClient.METHOD_POST, json_body)
	
	if error != OK:
		error_occurred.emit("HTTP Request failed to start")
		http_request.queue_free()

func _on_request_completed(result, response_code, headers, body, http_request):
	http_request.queue_free()
	
	if response_code != 200:
		error_occurred.emit("TTS API Error: " + str(response_code))
		return
	
	var stream = AudioStreamMP3.new()
	stream.data = body
	audio_received.emit(stream)
