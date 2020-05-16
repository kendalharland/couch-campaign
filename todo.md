# MVP Completion Criteria
1. One or more players can join the same game using an access code (no security)
2. Players can continue playing until the game disconnects or one is elected president.
   (deck re-shuffles on empty).
3. The UI looks like the mocks I drew (header functionality excluded).

# Start
* Writing tests as needed.
* Refactoring for neat-ness as-needed.

# MVP Todos
* [DONE] Auto-generate and assign leader names at start of game.
* [DONE] Auto-generate and re-assign leader names on state failure.
* [DONE] Randomly generate and order card decks (shuffle)
* [DONE] Regenerate the deck once it is empty.
* [DONE] Render the card according to the mockup I drew.
* Client|Server: Support joining game lobbies.
* Client: Create screen for player death.
* Client: Return to the main menu when the game is over.
* Client: Show a disconnect message and return to main menu when socket is closed.
* Client: Display time-in-office in footer.
* Server: Redesign the game state as func(ClientInput) GameOutput. Test that function

# Deferred
* [DONE] Come up with algorithm for determining who wins the election.
* Liar mechanic: Card that lets players insert a card into other player's decks.
* Quest mechanic: Card that lets players insert cards into their own decks + achievements. 

# Content
* Outlaw (marijuana, drugs)
* Stop illegal fracking
* Raise income taxes
* Raise taxes on local companies
* Endorse local politician running for mayor / small office

# House-Keeping
* Lobby server should reap dead game processes.
* Find a way to specify cards, deck order & effects outside the source code.

# Noteworthy
* Alternate text themes: 
  * Fonts
  * neucha
  * gamja flower
  * schoolbell
  * Patrick Hand SC
