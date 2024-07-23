package commander

func isGameCommand(cmd string) bool {
    return false
}

func GameCommand(cmd string) string {
    if ! isGameCommand(cmd) { return "Sorry, I didn't understand that..." }
    return ""
}
