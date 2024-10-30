package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"errors"
	"fmt"
	"html"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	json_parser "github.com/justinjest/gator/internal/config"
	"github.com/justinjest/gator/internal/database"

	_ "github.com/lib/pq"
)

type state struct {
	db  *database.Queries
	cfg *json_parser.Config
}
type command struct {
	name string
	args []string
}
type commands struct {
	method map[string]func(*state, command) error
}
type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}
type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	var feed RSSFeed
	var client http.Client
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return &feed, err
	}
	req.Header.Add("user-agent", "gator")
	res, err := client.Do(req)
	if err != nil {
		return &feed, err
	}
	defer res.Body.Close()
	b, err := io.ReadAll(res.Body)
	if err != nil {
		return &feed, err
	}
	err = xml.Unmarshal(b, &feed)
	if err != nil {
		return &feed, err
	}
	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)
	return &feed, nil
}
func (c *commands) register(name string, f func(*state, command) error) {
	c.method[name] = f
}
func (c *commands) run(s *state, cmd command) error {
	err := c.method[cmd.name](s, cmd)
	if err != nil {
		return err
	}
	return nil
}
func addfeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 2 {
		return errors.New("addfeed takes two commands, name and url")
	}
	now := time.Now()
	uuid1 := uuid.New().String()
	params := database.CreateFeedParams{
		ID:        uuid1,
		CreatedAt: now,
		UpdatedAt: now,
		Name:      cmd.args[0],
		Url:       cmd.args[1],
		UserID:    user.ID,
	}
	feed, err := s.db.CreateFeed(context.Background(), params)
	if err != nil {
		return err
	}
	now = time.Now()
	uuid2 := uuid.New().String()
	paramsFeedFollows := database.CreateFeedFollowParams{
		ID:        uuid2,
		CreatedAt: now,
		UpdatedAt: now,
		UserID:    user.ID,
		FeedID:    feed.ID,
	}
	_, err = s.db.CreateFeedFollow(context.Background(), paramsFeedFollows)
	if err != nil {
		return err
	}
	entry, err := s.db.GetFeed(context.Background(), feed.Name)
	if err != nil {
		return err
	}
	print("%v\n", entry.Name)
	return nil
}
func agg(s *state, cmd command) error {
	url := "https://www.wagslane.dev/index.xml"
	res, err := fetchFeed(context.Background(), url)
	if err != nil {
		return err
	}
	fmt.Printf("%v\n", res)
	return nil
}
func reset(s *state, cmd command) error {
	err := s.db.Reset(context.Background())
	if err != nil {
		return err
	}
	fmt.Printf("Reset complete, all accounts deleted\n")
	return nil
}
func getFeeds(s *state, cmd command) error {
	if len(cmd.args) != 0 {
		return errors.New("getFeeds accepts no arguments")
	}
	feeds, err := s.db.Pprint(context.Background())
	if err != nil {
		return err
	}
	for _, item := range feeds {
		fmt.Printf("%v\n", item)
	}
	return nil
}
func getUsers(s *state, cmd command) error {
	currentUser := s.cfg.Current_user_name
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return err
	}
	for _, user := range users {
		if user != currentUser {
			fmt.Printf("* %v\n", user)
		} else {
			fmt.Printf("* %v (current)\n", user)
		}
	}
	return nil
}
func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("user name requried for login function")
	}
	if len(cmd.args) >= 2 {
		return errors.New("login takes exactly one argument")
	}
	name := cmd.args[0]
	_, err := s.db.GetUser(context.Background(), name)
	if err != nil {
		return err
	}
	_, err = json_parser.SetUser(*s.cfg, name)
	if err != nil {
		return err
	}
	fmt.Printf("Username set to %v\n", name)
	return nil
}
func registerNewUser(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("user name requried for register function")
	}
	if len(cmd.args) >= 2 {
		return errors.New("register takes exactly one argument")
	}
	timeNow := time.Now()
	uuid := uuid.New().String()
	name := cmd.args[0]
	params := database.CreateUserParams{
		ID:        uuid,
		CreatedAt: timeNow,
		UpdatedAt: timeNow,
		Name:      name,
	}
	user, err := s.db.CreateUser(context.Background(), params)
	if err != nil {
		return err
	}
	s.cfg.Current_user_name = name
	fmt.Printf("User created %v", user)
	_, err = json_parser.SetUser(*s.cfg, name)
	if err != nil {
		return err
	}
	return nil
}
func follow(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 1 {
		return errors.New("follow takes exactly one argument, the url")
	}
	uuid := uuid.New().String()
	now := time.Now()
	feed, err := s.db.GetFeedByUrl(context.Background(), cmd.args[0])
	if err != nil {
		return err
	}
	params := database.CreateFeedFollowParams{
		ID:        uuid,
		CreatedAt: now,
		UpdatedAt: now,
		UserID:    user.ID,
		FeedID:    feed.ID,
	}
	_, err = s.db.CreateFeedFollow(context.Background(), params)
	if err != nil {
		return err
	}
	fmt.Printf("%v\n", user.Name)
	fmt.Printf("%v\n", feed.Name)
	return nil
}
func following(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 0 {
		return errors.New("following does not accept any arguments")
	}
	feeds, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return err
	}
	for _, data := range feeds {
		fmt.Printf("%v\n", data)
	}
	return nil
}
func main() {
	s, err := startUp()
	if err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}
	c := commands{
		method: make(map[string]func(*state, command) error),
	}
	c.register("login", handlerLogin)
	c.register("register", registerNewUser)
	c.register("reset", reset)
	c.register("users", getUsers)
	c.register("agg", agg)
	c.register("addfeed", middlewareLoggedIn(addfeed))
	c.register("feeds", getFeeds)
	c.register("follow", middlewareLoggedIn(follow))
	c.register("following", middlewareLoggedIn(following))
	if len(os.Args) < 2 {
		err = errors.New("too few cmdline arguments")
	}
	if err != nil {
		fmt.Printf("Error %v", err)
		os.Exit(1)
	}
	cmd := command{
		os.Args[1],
		os.Args[2:],
	}
	err = c.run(&s, cmd)
	if err != nil {
		fmt.Printf("Error %v\n", err)
		os.Exit(1)
	}
}
func startUp() (state, error) {
	config, err := json_parser.Read()
	if err != nil {
		return state{}, errors.New("unable to read json")
	}
	var s state
	s.cfg = &config
	db, err := sql.Open("postgres", s.cfg.Db_url)
	if err != nil {
		return state{}, errors.New("unable to open postgres db")
	}
	dbQueries := database.New(db)
	s.db = dbQueries
	return s, nil
}
func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		user, err := s.db.GetUser(context.Background(), s.cfg.Current_user_name)
		if err != nil {
			return err
		}
		return handler(s, cmd, user)
	}
}
