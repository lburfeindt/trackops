on run arguments
	set playlistPath to item 1 of arguments
	tell application "TrackOps" to activate
	tell application "System Events"
		tell process "TrackOps"
			repeat 30 times
				if exists window "TrackOps" then exit repeat
				delay 0.2
			end repeat

			set openButton to missing value
			repeat with element in entire contents of window "TrackOps"
				if role of element is "AXButton" and name of element is "Format Playlist" then
					set openButton to element
					exit repeat
				end if
			end repeat
			if openButton is missing value then error "Format Playlist button is missing"
			click openButton

			keystroke "g" using {command down, shift down}
			keystroke playlistPath
			key code 36
			repeat 30 times
				if exists button "Open" of sheet 1 of window "TrackOps" then exit repeat
				delay 0.2
			end repeat
			click button "Open" of sheet 1 of window "TrackOps"

			set expectedTitles to {"Example Artist - Example Track", "Artist One, Artist Two - Café Song, Club Mix"}
			repeat 30 times
				set displayedTitles to {}
				repeat with element in entire contents of window "TrackOps"
					try
						set displayedValue to value of element as text
						if expectedTitles contains displayedValue then set end of displayedTitles to displayedValue
					end try
				end repeat
				if displayedTitles is equal to expectedTitles then return
				delay 0.2
			end repeat
			error "Playlist titles were not visible in the TrackOps window"
		end tell
	end tell
end run
