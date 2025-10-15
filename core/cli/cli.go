// package cli provides a lightweight, zero-dependency CLI framework.
package cli

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"chef/core/twx"
)

type Executable interface {
	Exec(ctx context.Context, args ...string) error
	GetName() string
	Print() string
}

type App struct {
	Menu *Menu
}

func NewApp(name, overview string) *App {
	menu := &Menu{
		Name:     name,
		Overview: overview,
		Execs:    make([]Executable, 0),
	}

	app := &App{menu}

	return app
}

func (a *App) Run(args ...string) error {
	ctx := context.Background()

	return a.Menu.Exec(ctx, args...)
}

// ====================================================================
// CONVENIENCE PROXY METHODS

func (a *App) AddFunc(name, desc string, f Func) {
	a.Menu.AddFunc(name, desc, f)
}

func (a *App) AddGroup(name string) {
	a.Menu.AddGroup(name)
}

func (a *App) AddMenu(name, desc, overview string) *Menu {
	return a.Menu.AddMenu(name, desc, overview)
}

// ====================================================================
// FUNC

type Func func(ctx context.Context, args ...string) error

type fn struct {
	f    Func
	name string
	desc string
}

func (f fn) Exec(ctx context.Context, args ...string) error {
	return f.f(ctx, args...)
}

func (f fn) Print() string {
	return fmt.Sprintf("\t%s\t%s\n", f.name, f.desc)
}

func (f fn) GetName() string {
	return f.name
}

// ====================================================================
// GROUP

type Group string

// Exec is a dummy method implemented to comply with the Executable interface
func (g Group) Exec(ctx context.Context, args ...string) error {
	return nil
}

func (g Group) Print() string {
	return fmt.Sprintf("\n%s:\n", g)
}

// GetName is a dummy method implemented to comply with the Executable interface
func (g Group) GetName() string {
	return ""
}

// ====================================================================
// MENU

type Menu struct {
	Name     string
	Desc     string
	Overview string
	Execs    []Executable

	// this is just to decide whether a space must be printed
	// below the overview or not.
	firstEntryIsGroup bool
}

func (m *Menu) AddFunc(name, desc string, f Func) {
	validateName(name)

	ff := fn{
		f:    f,
		name: name,
		desc: desc,
	}

	m.Execs = append(m.Execs, ff)
}

func (m *Menu) AddGroup(name string) {
	if len(m.Execs) == 0 {
		m.firstEntryIsGroup = true
	}

	m.Execs = append(m.Execs, Group(name))
}

func (m *Menu) AddMenu(name, desc, overview string) *Menu {
	mm := &Menu{
		Name:     name,
		Desc:     desc,
		Overview: overview,
		Execs:    make([]Executable, 0),
	}

	m.Execs = append(m.Execs, mm)

	return mm
}

func (m Menu) Exec(ctx context.Context, args ...string) error {
	// if there are no further args, we need to print the menu
	if len(args) == 0 {
		m.printMenu()
		return nil
	}

	for _, exec := range m.Execs {
		if exec.GetName() == args[0] {
			return exec.Exec(ctx, args[1:]...)
		}
	}

	return fmt.Errorf("invalid argument %q", args[0])
}

func (m *Menu) Print() string {
	return fmt.Sprintf("\t%s\t%s\n", m.Name, m.Desc)
}

func (m *Menu) GetName() string {
	return m.Name
}

func (m *Menu) printMenu() {
	var builder strings.Builder
	tw := twx.NewWriter(&builder)

	fmt.Fprintf(tw, "%s\n", m.Overview)
	if !m.firstEntryIsGroup {
		fmt.Fprint(tw, "\n")
	}

	for _, e := range m.Execs {
		fmt.Fprint(tw, e.Print())
	}

	tw.Flush()

	fmt.Print(builder.String())
}

// ====================================================================
// INVARIANTS

func initErr(s string) error {
	return fmt.Errorf("CLI Init Error: %s", s)
}

func validateName(name string) {
	if strings.Contains(name, " ") {
		panic(initErr(fmt.Sprintf("invalid arg name %q, spaces not allowed", name)))
	}
}

// ====================================================================
// Input

func InteractiveStr(question string, answer *string) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print(question)
	val, _ := reader.ReadString('\n')
	val = strings.TrimSuffix(val, "\n")
	*answer = val
}

func InteractiveYesNo(question string, answer *bool) {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print(question)
		val, _ := reader.ReadString('\n')
		val = strings.TrimSuffix(val, "\n")

		if val == "Y" || val == "y" || val == "yes" || val == "Yes" || val == "YES" {
			*answer = true
			break
		}

		if val == "N" || val == "n" || val == "no" || val == "No" || val == "NO" {
			*answer = true
			break
		}
	}

}
