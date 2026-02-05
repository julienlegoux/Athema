extends Node2D

# Animation names - ADJUST these to match your actual Spine export
const ANIM_IDLE = "idle"
const ANIM_SPEAK = "speak"
const ANIM_LISTEN = "listen"

@onready var spine_sprite = $SpineSprite

func _ready():
	# If we are in mock mode and no sprite exists, we warn but don't crash
	if not has_node("SpineSprite"):
		print("SpineCharacter: Warning - 'SpineSprite' node not found. Animations will be log-only.")

func play_idle():
	_play_animation(ANIM_IDLE, true)

func play_speak():
	_play_animation(ANIM_SPEAK, true)

func play_listen():
	_play_animation(ANIM_LISTEN, true)

func _play_animation(anim_name: String, loop: bool):
	if has_node("SpineSprite"):
		# This syntax depends on the specific Godot-Spine runtime API
		# Usually it's get_animation_state().set_animation(anim_name, loop, 0)
		var sprite = get_node("SpineSprite")
		if sprite.has_method("get_animation_state"):
			sprite.get_animation_state().set_animation(anim_name, loop, 0)
		else:
			print("SpineCharacter: SpineSprite found but API mismatch.")
	else:
		print("SpineCharacter (Mock): Playing animation -> ", anim_name)
