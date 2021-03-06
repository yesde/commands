// Copyright 2013 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package commands

import (
	"testing"

	"github.com/limetext/backend"
	"github.com/limetext/text"
)

type caseTest struct {
	inRegion []text.Region
	in       string
	exp      string
}

func runCaseTest(command string, testsuite *[]caseTest, t *testing.T) {
	ed := backend.GetEditor()
	w := ed.NewWindow()
	defer w.Close()

	for i, test := range *testsuite {
		v := w.NewFile()
		defer func() {
			v.SetScratch(true)
			v.Close()
		}()

		e := v.BeginEdit()
		v.Insert(e, 0, test.in)
		v.EndEdit(e)

		v.Sel().Clear()
		if test.inRegion != nil {
			for _, r := range test.inRegion {
				v.Sel().Add(r)
			}
		}
		ed.CommandHandler().RunTextCommand(v, command, nil)
		sr := v.Substr(text.Region{0, v.Size()})
		if sr != test.exp {
			t.Errorf("%s test %d failed: %v, %+v", command, i, sr, test)
		}
	}
}

func TestTitleCase(t *testing.T) {
	tests := []caseTest{
		/*single selection*/
		{
			// Please note the bizarre  capitalization of the first L in he'Ll...  This is due to a bug in go's strings
			// library.  I'm going to try to get them to fix it...  If not, maybe we'll have
			// to write our own Title Casing function.
			[]text.Region{{24, 51}},

			"Give a man a match, and he'll be warm for a minute, but set him on fire, and he'll be warm for the rest of his life.",
			"Give a man a match, and He'Ll Be Warm For A Minute, but set him on fire, and he'll be warm for the rest of his life.",
		},
		/*multiple selection*/
		{
			[]text.Region{{0, 17}, {52, 71}},

			"Give a man a match, and he'll be warm for a minute, but set him on fire, and he'll be warm for the rest of his life.",
			"Give A Man A Match, and he'll be warm for a minute, But Set Him On Fire, and he'll be warm for the rest of his life.",
		},
		/*no selection*/
		{
			nil,

			"Give a man a match, and he'll be warm for a minute, but set him on fire, and he'll be warm for the rest of his life.",
			"Give a man a match, and he'll be warm for a minute, but set him on fire, and he'll be warm for the rest of his life.",
		},
		/*unicode*/
		{
			[]text.Region{{0, 12}},

			"ничего себе!",
			"Ничего Себе!",
		},
		/*asian characters*/
		{
			[]text.Region{{0, 9}},

			"千里之行﹐始于足下",
			"千里之行﹐始于足下",
		},
	}

	runCaseTest("title_case", &tests, t)
}

func TestSwapCase(t *testing.T) {
	tests := []caseTest{
		{
			[]text.Region{{0, 0}},

			"",
			"",
		},
		{
			[]text.Region{{0, 13}},

			"Hello, World!",
			"hELLO, wORLD!",
		},
		{
			[]text.Region{{0, 11}},

			"ПрИвЕт, МиР",
			"пРиВеТ, мИр",
		},
	}

	runCaseTest("swap_case", &tests, t)
}

func TestUpperCase(t *testing.T) {
	tests := []caseTest{
		/*single selection*/
		{
			[]text.Region{{0, 76}},

			"Try not to become a man of success, but rather try to become a man of value.",
			"TRY NOT TO BECOME A MAN OF SUCCESS, BUT RATHER TRY TO BECOME A MAN OF VALUE.",
		},
		/*multiple selection*/
		{
			[]text.Region{{0, 20}, {74, 76}},

			"Try not to become a man of success, but rather try to become a man of value.",
			"TRY NOT TO BECOME A man of success, but rather try to become a man of valuE.",
		},
		/*no selection*/
		{
			nil,

			"Try not to become a man of success, but rather try to become a man of value.",
			"Try not to become a man of success, but rather try to become a man of value.",
		},
		/*unicode*/
		{
			[]text.Region{{0, 74}},

			"чем больше законов и постановлений, тем больше разбойников и преступлений!",
			"ЧЕМ БОЛЬШЕ ЗАКОНОВ И ПОСТАНОВЛЕНИЙ, ТЕМ БОЛЬШЕ РАЗБОЙНИКОВ И ПРЕСТУПЛЕНИЙ!",
		},
		/*asian characters*/
		{
			[]text.Region{{0, 9}},

			"千里之行﹐始于足下",
			"千里之行﹐始于足下",
		},
	}

	runCaseTest("upper_case", &tests, t)
}

func TestLowerCase(t *testing.T) {
	tests := []caseTest{
		/*single selection*/
		{
			[]text.Region{{0, 76}},

			"TRY NOT TO BECOME A MAN OF SUCCESS, BUT RATHER TRY TO BECOME A MAN OF VALUE.",
			"try not to become a man of success, but rather try to become a man of value.",
		},
		/*multiple selection*/
		{
			[]text.Region{{0, 20}, {74, 76}},

			"TRY NOT TO BECOME A MAN OF SUCCESS, BUT RATHER TRY TO BECOME A MAN OF VALUE.",
			"try not to become a MAN OF SUCCESS, BUT RATHER TRY TO BECOME A MAN OF VALUe.",
		},
		/*no selection*/
		{
			nil,

			"Try not to become a man of success, but rather try to become a man of value.",
			"Try not to become a man of success, but rather try to become a man of value.",
		},
		/*unicode*/
		{
			[]text.Region{{0, 74}},

			"ЧЕМ БОЛЬШЕ ЗАКОНОВ И ПОСТАНОВЛЕНИЙ, ТЕМ БОЛЬШЕ РАЗБОЙНИКОВ И ПРЕСТУПЛЕНИЙ!",
			"чем больше законов и постановлений, тем больше разбойников и преступлений!",
		},
		/*asian characters*/
		{
			[]text.Region{{0, 9}},

			"千里之行﹐始于足下",
			"千里之行﹐始于足下",
		},
	}

	runCaseTest("lower_case", &tests, t)
}
