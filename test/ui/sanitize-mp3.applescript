on run arguments
	-- Given
	set selectionDirectory to item 1 of arguments
	tell application "TrackOps" to activate

	-- When
	my openPicker()
	tell application "System Events" to tell process "TrackOps"
		key code 53
	end tell

	-- Then
	my waitForText("Choose MP3 files to remove the (Original Mix) title suffix.")

	-- When
	my chooseFiles(selectionDirectory)

	-- Then
	my waitForText("Updated: 1 · Unchanged: 1 · Failed: 1")
	my waitForTextContaining("broken.mp3")

	-- When
	my chooseFiles(selectionDirectory)

	-- Then
	my waitForText("Updated: 0 · Unchanged: 2 · Failed: 1")
end run

on findButton(buttonName)
	tell application "System Events" to tell process "TrackOps"
		repeat 30 times
			set windowElements to get entire contents of window "TrackOps"
			repeat with element in windowElements
				if role of element is "AXButton" and name of element is buttonName then return element
			end repeat
			delay 0.2
		end repeat
	end tell
	error buttonName & " button is missing"
end findButton

on openPicker()
	tell application "System Events" to tell process "TrackOps"
		click my findButton("ID3 Sanitizer")
		repeat 30 times
			if exists sheet 1 of window "TrackOps" then return
			delay 0.2
		end repeat
	end tell
	error "Native MP3 picker did not open"
end openPicker

on chooseFiles(selectionDirectory)
	my openPicker()
	tell application "System Events" to tell process "TrackOps"
		set frontmost to true
		keystroke "g" using {command down, shift down}
		repeat 30 times
			if exists sheet 1 of sheet 1 of window "TrackOps" then exit repeat
			delay 0.2
		end repeat
		set value of text field 1 of sheet 1 of sheet 1 of window "TrackOps" to selectionDirectory
		key code 36
		repeat 30 times
			if not exists sheet 1 of sheet 1 of window "TrackOps" then exit repeat
			delay 0.2
		end repeat
		keystroke "a" using command down
		click button "Open" of splitter group 1 of sheet 1 of window "TrackOps"
	end tell
end chooseFiles

on waitForText(expectedText)
	repeat 50 times
		if my hasText(expectedText) then return
		delay 0.2
	end repeat
	error expectedText & " was not visible"
end waitForText

on waitForTextContaining(expectedText)
	repeat 50 times
		tell application "System Events" to tell process "TrackOps"
			set windowElements to get entire contents of window "TrackOps"
			repeat with element in windowElements
				try
					if (value of element as text) contains expectedText then return
				end try
			end repeat
		end tell
		delay 0.2
	end repeat
	error expectedText & " was not visible"
end waitForTextContaining

on hasText(expectedText)
	tell application "System Events" to tell process "TrackOps"
		set windowElements to get entire contents of window "TrackOps"
		repeat with element in windowElements
			try
				if value of element as text is expectedText then return true
			end try
		end repeat
	end tell
	return false
end hasText
