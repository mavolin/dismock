# dismock

[![GitHub Workflow Status (branch)](https://img.shields.io/github/workflow/status/mavolin/dismock/Test/v2)](https://github.com/mavolin/dismock/actions?query=workflow%3ATest+branch%3Av2+)
[![Test Coverage](https://codecov.io/gh/mavolin/dismock/branch/v2/graph/badge.svg)](https://codecov.io/gh/mavolin/dismock/branch/v2)
[![Go Report Card](https://goreportcard.com/badge/github.com/mavolin/dismock)](https://goreportcard.com/report/github.com/mavolin/dismock)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/mavolin/dismock/v2)](https://pkg.go.dev/github.com/mavolin/dismock/v2)
[![License](https://img.shields.io/github/license/mavolin/dismock)](https://github.com/mavolin/dismock/blob/v2/LICENSE)

-----

Dismock is a library that aims to make mocking Discord's API requests as easy as winking.
No more huge integration tests that require a bot on some private server with little to no debug information.

Dismock is not limited to a specific Discord library, although it uses [arikawa](https://github.com/diamondburned/arikawa) as a foundation for its datatypes.

## Getting Started

#### Basic Testing

Creating a mock is done, by calling the respective mock method, that belongs to the API request you made in your code.
Below is a basic example of a ping command and it's unit test.

```go
func (b *Bot) Ping(e *gateway.MessageCreateEvent) (error) {
    _, err := b.Ctx.SendText(e.ChannelID, "🏓")
    if err != nil {
        return err
    }

    _, err := b.Ctx.SendText(e.ChannelID, "Pong!")
    
    return err
}
```

```go
func TestBot_Ping(t *testing.T) {
    m, s := dismock.NewState(t)
    // if you want to use a Session and no State write:
    // m, s := dismock.NewSession(t)

    var channelID discord.ChannelID = 123

    m.SendText(discord.Message{
        // from the doc of Mocker.SendText we know, that ChannelID and Content
    	// are required fields, all other fields that aren't used by our function 
    	// don't need to be filled
        ChannelID: channelID,
        Content: "🏓",
    })

    // make sure that API calls on the same endpoint with the same method are
    // added in the correct order
    m.SendText(discord.Message{
        ChannelID: channelID,
        Content: "Pong!"
    })

    b := NewBot(s)

    b.Ping(&gateway.MessageCreateEvent{
        Message: discord.Message{
            ChannelID: channelID,
        }
    })

    // at the end of every test Mocker.Eval must be called, to check for
    // uninvoked handlers and close the mock server
    m.Eval()
}
```

#### Advanced Testing

Now imagine a bit more complicated test, that has multiple sub-tests:

```go
func (b *Bot) Ping(e *gateway.MessageCreateEvent) (error) {
    _, err := b.Ctx.SendText(e.ChannelID, "🏓")
    if err != nil {
        return err
    }

    _, err := b.Ctx.SendText(e.ChannelID, e.Author.Mention()+" Pong!")
    
    return err
}
```

```go
func TestBot_Ping(t *testing.T) {
    m := dismock.New(t)

    var channelID discord.ChannelID = 123

    m.SendText(discord.Message{
        ChannelID: channelID,
        Content: "🏓",
    })
    
   t.Run("test1", func(t *testing.T) {
        // If you have multiple tests that have the same basic API requests,
        // you can create a mocker, add those API calls, and create a clone
        // of the mocker in every sub-test you have.
        // Cloned mockers have a copy of it's parents request, but run their
        // own mock server and have a dedicated Session/State.
        m, s := m.CloneState(t)

        ...
    })

    t.Run("test2", func(t *testing.T) {
        m, s := m.CloneState(t)
        
        ...
    })
}
```

### Using a Different Discord Library

Because the mocking is done on a network level, you are not limited to arikawa for mocking.
When creating a `Mocker` just use `dismock.New` and use the `http.Client` of the mocker as client for a library of your choice.

```go
m := dismock.New(t)

s, _ := discordgo.New("Bot abc") // the token doesn't have to be valid
s.Client = m.Client
```

That's it!

### Meta Requests

Besides regular calls to the API you can also mock requests for metadata, i.e. images such as guild icons (`Mocker.GuildIcon`).
In order for this to work, you need to use the `http.Client` found in the `Mocker` struct, so that the mock server will be called, instead of Discord.
