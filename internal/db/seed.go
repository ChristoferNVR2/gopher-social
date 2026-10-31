package db

import (
	"context"
	"fmt"
	"log"
	"math/rand"

	"github.com/ChristoferNVR2/gopher-social/internal/models"
	"github.com/ChristoferNVR2/gopher-social/internal/store"
)

var usernames = []string{
	"alex", "ava", "ben", "chris", "claire", "daniel", "ella", "emily", "ethan", "felix",
	"grace", "harper", "henry", "isla", "jack", "james", "jasmine", "julia", "kai", "leo",
	"liam", "lily", "lucas", "lucy", "mason", "maya", "mia", "nathan", "noah", "nora",
	"oliver", "olivia", "oscar", "owen", "peter", "quinn", "rachel", "ruby", "sam", "sarah",
	"sebastian", "sofia", "theo", "thomas", "victor", "violet", "will", "willow", "wyatt", "zoe",
}

var titles = []string{
	"Getting Started with Go", "Building Your First REST API", "Understanding Go Interfaces",
	"Error Handling Best Practices", "Concurrency with Goroutines", "Working with Channels in Go",
	"Writing Clean and Maintainable Code", "Introduction to Database Design", "Testing Go Applications",
	"Deploying a Go Web Service", "Authentication and Authorization", "Structuring a Go Project",
	"Useful Git Workflows for Developers", "Improving Application Performance", "Logging and Monitoring Basics",
	"Designing Better APIs", "Dependency Management in Go", "Working with JSON Data",
	"Common Backend Development Patterns", "Lessons Learned from Building a Social App",
}

var contents = []string{
	"Go is easy to learn when you start with small programs, understand the basic syntax, and practice building useful command-line tools.",
	"A REST API exposes resources through predictable HTTP endpoints, making it easier for clients to create, read, update, and delete data.",
	"Interfaces in Go describe behavior rather than implementation, allowing components to remain flexible, testable, and loosely coupled.",
	"Clear errors, useful context, and deliberate handling of failures make Go applications easier to debug and safer to operate.",
	"Goroutines make concurrent work lightweight, but shared state should be protected and every goroutine should have a clear lifecycle.",
	"Channels provide a structured way for goroutines to communicate, coordinate work, and avoid unnecessary shared-memory complexity.",
	"Readable names, small functions, and focused packages help code remain understandable as an application grows.",
	"A thoughtful database design starts with clear entities, appropriate relationships, useful constraints, and indexes based on real queries.",
	"Table-driven tests and meaningful assertions make Go applications more reliable while keeping test suites simple to extend.",
	"Reliable deployment requires a repeatable build process, configuration through the environment, health checks, and sensible rollback plans.",
	"Authentication confirms who a user is, while authorization determines what that user is allowed to do; both need careful security controls.",
	"A consistent project structure separates responsibilities and makes it easier for developers to find, test, and change code.",
	"Small commits, descriptive branch names, and regular rebasing can make collaboration smoother and simplify code reviews.",
	"Measure before optimizing: profiling and tracing reveal the real bottlenecks so performance work produces meaningful improvements.",
	"Useful logs, metrics, and alerts provide the visibility needed to understand production behavior and respond to incidents quickly.",
	"Good APIs use consistent naming, status codes, validation, pagination, and documentation so clients can depend on predictable behavior.",
	"Go modules make dependency versions explicit and reproducible, while regular updates help address bugs and security issues.",
	"JSON is convenient for exchanging data, but explicit structs and validation keep serialization behavior predictable and safe.",
	"Handlers, services, repositories, and background workers each have distinct responsibilities that can be composed into a maintainable backend.",
	"Building a social app teaches valuable lessons about data modeling, user experience, moderation, scalability, and the importance of feedback.",
}

var tags = []string{
	"golang", "programming", "backend", "web-development", "api", "databases", "testing", "devops", "cloud", "security",
	"performance", "architecture", "git", "javascript", "frontend", "open-source", "career", "tutorial", "projects", "community",
}

var comments = []string{
	"Great post! I learned a lot.",
	"Thanks for sharing your insights.",
	"This was really helpful, especially the part about error handling.",
	"I appreciate the detailed explanation of interfaces.",
	"Can you provide more examples of using channels?",
	"Your tips on testing are very useful.",
	"I'm looking forward to your next post on deployment strategies.",
	"This article clarified a lot of my doubts about Go concurrency.",
	"Excellent write-up! Keep up the good work.",
	"I found the section on database design particularly enlightening.",
}

func Seed(store store.Storage) {
	ctx := context.Background()

	users := generateUsers(100)
	for _, user := range users {
		if err := store.Users.Create(ctx, user); err != nil {
			log.Println("failed to create user: ", err)
			return
		}
	}

	posts := generatePosts(200, users)
	for _, post := range posts {
		if err := store.Posts.Create(ctx, post); err != nil {
			log.Println("failed to create post: ", err)
			return
		}
	}

	cms := generateComments(500, posts, users)
	for _, comment := range cms {
		if err := store.Comments.Create(ctx, comment); err != nil {
			log.Println("failed to create comment: ", err)
			return
		}
	}

	log.Println("Database seeded successfully with users, posts, and comments.")
}

func generateUsers(num int) []*models.User {
	users := make([]*models.User, num)

	for i := 0; i < num; i++ {
		users[i] = &models.User{
			Username: usernames[i%len(usernames)] + fmt.Sprintf("%d", i),
			Email:    usernames[i%len(usernames)] + fmt.Sprintf("%d", i) + "@example.com",
			Password: "123123",
		}
	}

	return users
}

func generatePosts(num int, users []*models.User) []*models.Post {
	posts := make([]*models.Post, num)

	for i := 0; i < num; i++ {
		user := users[rand.Intn(len(users))]

		posts[i] = &models.Post{
			Title:   titles[rand.Intn(len(titles))],
			Content: contents[rand.Intn(len(contents))],
			UserID:  user.ID,
			Tags: []string{
				tags[rand.Intn(len(tags))],
				tags[rand.Intn(len(tags))],
			},
		}
	}

	return posts
}

func generateComments(num int, posts []*models.Post, users []*models.User) []*models.Comment {
	cms := make([]*models.Comment, num)

	for i := 0; i < num; i++ {
		cms[i] = &models.Comment{
			PostID:  posts[rand.Intn(len(posts))].ID,
			UserID:  users[rand.Intn(len(users))].ID,
			Content: comments[rand.Intn(len(comments))],
		}
	}

	return cms
}
