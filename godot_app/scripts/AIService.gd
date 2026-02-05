extends Node
class_name AIService

signal response_received(text: String)
signal error_occurred(message: String)

# Anthropic API Config
const API_URL = "https://api.anthropic.com/v1/messages"
const MODEL = "claude-3-haiku-20240307"

func send_message(user_text: String):
	if Global.use_mocks:
		_mock_send_message(user_text)
	else:
		_real_send_message(user_text)

func _mock_send_message(user_text: String):
	print("AIService (Mock): Sending message -> ", user_text)
	# Simulate network delay
	await get_tree().create_timer(1.5).timeout
	
	var mock_responses = [
		"I am Athema, your digital companion.",
		"That is an interesting perspective.",
		"I am listening diligently.",
		"Could you elaborate on that?",
		"My internal systems are functioning perfectly."
	]
	var response = mock_responses.pick_random()
	print("AIService (Mock): Received response -> ", response)
	response_received.emit(response)

func _real_send_message(user_text: String):
	if Global.claude_api_key.is_empty():
		error_occurred.emit("Missing Claude API Key")
		return

	var http_request = HTTPRequest.new()
	add_child(http_request)
	http_request.request_completed.connect(_on_request_completed.bind(http_request))

	var headers = [
		"x-api-key: " + Global.claude_api_key,
		"anthropic-version: 2023-06-01",
		"content-type: application/json"
	]
	
	var body = {
		"model": MODEL,
		"max_tokens": 150,
		"messages": [
			{"role": "user", "content": user_text}
		]
	}
	
	var json_body = JSON.stringify(body)
	var error = http_request.request(API_URL, headers, HTTPClient.METHOD_POST, json_body)
	
	if error != OK:
		error_occurred.emit("HTTP Request failed to start")
		http_request.queue_free()

func _on_request_completed(result, response_code, headers, body, http_request):
	http_request.queue_free()
	
	if response_code != 200:
		error_occurred.emit("API Error: " + str(response_code))
		print("Body: ", body.get_string_from_utf8())
		return
	
	var json = JSON.new()
	json.parse(body.get_string_from_utf8())
	var response_data = json.get_data()
	
	if response_data.has("content") and response_data["content"].size() > 0:
		var text = response_data["content"][0]["text"]
		response_received.emit(text)
	else:
		error_occurred.emit("Unexpected API response format")
