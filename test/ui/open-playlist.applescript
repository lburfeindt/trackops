on run arguments
	-- Given
	set playlistPath to item 1 of arguments
	set playlistName to item 2 of arguments
	tell application "TrackOps" to activate
	tell application "System Events"
		tell process "TrackOps"
			repeat 30 times
				if exists window "TrackOps" then exit repeat
				delay 0.2
			end repeat

			set openButton to missing value
			repeat 30 times
				set windowElements to get entire contents of window "TrackOps"
				repeat with element in windowElements
					if role of element is "AXButton" and name of element is "Format Playlist" then
						set openButton to element
						exit repeat
					end if
				end repeat
				if openButton is not missing value then exit repeat
				delay 0.2
			end repeat
			if openButton is missing value then error "Format Playlist button is missing"

			-- When
			click openButton

			repeat 30 times
				if exists sheet 1 of window "TrackOps" then exit repeat
				delay 0.2
			end repeat
			set frontmost to true
			keystroke "g" using {command down, shift down}
			repeat 30 times
				if exists sheet 1 of sheet 1 of window "TrackOps" then exit repeat
				delay 0.2
			end repeat
			set value of text field 1 of sheet 1 of sheet 1 of window "TrackOps" to playlistPath
			key code 36

			repeat 30 times
				try
					if enabled of button "Open" of splitter group 1 of sheet 1 of window "TrackOps" then exit repeat
				end try
				delay 0.2
			end repeat
			if not enabled of button "Open" of splitter group 1 of sheet 1 of window "TrackOps" then error "Playlist was not selected in the native picker"
			click button "Open" of splitter group 1 of sheet 1 of window "TrackOps"

			-- Then
			repeat 30 times
				if not exists sheet 1 of window "TrackOps" then exit repeat
				delay 0.2
			end repeat
			if exists sheet 1 of window "TrackOps" then error "Native playlist picker did not close"

			set expectedTitles to {"Example Artist - Example Track", "Artist One, Artist Two - Café Song, Club Mix"}
			repeat 30 times
				set displayedTitles to {}
				set displayedValues to {}
				set windowElements to get entire contents of window "TrackOps"
				repeat with element in windowElements
					try
						set displayedValue to value of element as text
						set end of displayedValues to displayedValue
						if role of element is "AXStaticText" and expectedTitles contains displayedValue then set end of displayedTitles to displayedValue
					end try
				end repeat
				if displayedValues contains playlistName and displayedTitles is equal to expectedTitles then return
				delay 0.2
			end repeat
			error "Playlist titles were not visible in the TrackOps window"
		end tell
	end tell
end run
