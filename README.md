<div align="center">
<h1>dismock</h1>
    
[![GitHub Workflow Status (branch)](https://img.shields.io/github/workflow/status/mavolin/dismock/Test/v2)](https://github.com/mavolin/dismock/actions?query=workflow%3ATest+branch%3Av2+)
[![Test Coverage](https://codecov.io/gh/mavolin/dismock/branch/v2/graph/badge.svg)](https://codecov.io/gh/mavolin/dismock/branch/v2)
[![Go Report Card](https://goreportcard.com/badge/github.com/mavolin/dismock)](https://goreportcard.com/report/github.com/mavolin/dismock)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/mavolin/dismock/v2)](https://pkg.go.dev/github.com/mavolin/dismock/v2)
[![License](https://img.shields.io/github/license/mavolin/dismock)](https://github.com/mavolin/dismock/blob/v2/LICENSE)
</div>

---

Dismock is a library that aims to make mocking Discord's API requests as easy as winking.
No more huge integration tests that require a bot on some private server with little to no debug information.

Although, dismock uses [arikawa](https://github.com/diamondburned/arikawa) as a foundation for its data types, it isn't limited to a specific discord library.

## Getting Started

#### Basic Testing

You can create a mock by calling the method that corresponds to the API request you made in your code.
Below is a simple example of a ping command, and it's unit test.

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
    // you can also mock a Session by using dismock.NewSession(t), or dismock.New(t) 
    // to only create a Mocker
    m, s := dismock.NewState(t)

    var channelID discord.ChannelID = 123

    m.SendText(discord.Message{
        // the doc of every mock specifies what fields are required, all other
        // fields not relevant to your test can be omitted
        ChannelID: channelID,
        Content: "🏓",
    })

    // Mocks should be added in the same order their calls are made.
    // However, this order will only be enforced on calls to the same endpoint
    // using the same http method.
    m.SendText(discord.Message{
        ChannelID: channelID,
        Content: "Pong!"
    })

    b := NewBot(s)
    b.Ping(&gateway.MessageCreateEvent{
        Message: discord.Message{ChannelID: channelID}
    })
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
    // no need to call m.Eval() as we'll only use the mocker for cloning anyway

    var channelID discord.ChannelID = 123

    m.SendText(discord.Message{
        ChannelID: channelID,
        Content: "🏓",
    })

    t.Run("test1", func(t *testing.T) {
        // If you have multiple tests that make the same requests, you can
        // create a mocker, and add those API calls.
        // Afterwards, you can create a clone of the mocker in every sub-test 
        // you have.
        // Cloned mockers have a copy of their parent's request, but run their
        // own mock server and have a dedicated Session/State.
        m, s := m.CloneState(t)
        defer m.Eval()

        ...
    })

    t.Run("test2", func(t *testing.T) {
        m, s := m.CloneState(t)
        defer m.Eval()

        ...
    })
}
```

### Using a Different Discord Library

Since mocking is done on a network level, you are free to chose whatever discord library you want.
Simply use `dismock.New` when creating a mocker, replace the `http.Client` of your library of choice with `mocker.Client`, and disable the state.

Below is an example of using dismock with [discordgo](https://github.com/bwmarrin/discordgo).
```go
m := dismock.New(t)

s, _ := discordgo.New("Bot abc") // the token doesn't have to be valid
s.StateEnabled = false
s.Client = m.Client
```

### Meta Requests

Besides regular calls to the API, you can also mock requests for metadata, i.e. images such as guild icons (`Mocker.GuildIcon`).
In order for this to work you need to use the `http.Client` found in the `Mocker` struct, so that the mock server will be called instead of Discord.
