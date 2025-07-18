// Code generated by templ - DO NOT EDIT.

// templ: version: v0.3.906
package views

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

import (
	"fmt"
	"github.com/rbc33/gocms/common"

	components "github.com/rbc33/gocms/views/components"
)

func makeDarkModeButton() templ.Component {
	return templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
		templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context
		if templ_7745c5c3_CtxErr := ctx.Err(); templ_7745c5c3_CtxErr != nil {
			return templ_7745c5c3_CtxErr
		}
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templruntime.GetBuffer(templ_7745c5c3_W)
		if !templ_7745c5c3_IsBuffer {
			defer func() {
				templ_7745c5c3_BufErr := templruntime.ReleaseBuffer(templ_7745c5c3_Buffer)
				if templ_7745c5c3_Err == nil {
					templ_7745c5c3_Err = templ_7745c5c3_BufErr
				}
			}()
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 1, "<button id=\"theme-toggle\" class=\"navbar navbar-element inline-block focus:outline-none ml-4 transition duration-500\"><svg id=\"light-icon\" class=\"w-[30px] h-[30px] dark:block hidden\"><circle cx=\"15\" cy=\"15\" r=\"6\" fill=\"currentColor\"></circle> <line id=\"ray\" stroke=\"currentColor\" storke-width=\"2\" stroke-linecap=\"round\" x1=\"15\" y1=\"1\" x2=\"15\" y2=\"4\"></line> <use href=\"#ray\" transform=\"rotate(45 15 15)\"></use> <use href=\"#ray\" transform=\"rotate(90 15 15)\"></use> <use href=\"#ray\" transform=\"rotate(135 15 15)\"></use> <use href=\"#ray\" transform=\"rotate(180 15 15)\"></use> <use href=\"#ray\" transform=\"rotate(225 15 15)\"></use> <use href=\"#ray\" transform=\"rotate(270 15 15)\"></use> <use href=\"#ray\" transform=\"rotate(315 15 15)\"></use></svg> <svg id=\"dark-icon\" class=\"dark:hidden block w-[30px] h-[30px] rotate-[315deg]\"><path fill=\"currentColor\" d=\"M 23, 5 A 12 12 0 1 0 23, 25  A 12 12 0 0 1 23, 5\"></path></svg></button>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return nil
	})
}

func MakeNavBar(links []common.Link, dropdowns map[string][]common.Link) templ.Component {
	return templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
		templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context
		if templ_7745c5c3_CtxErr := ctx.Err(); templ_7745c5c3_CtxErr != nil {
			return templ_7745c5c3_CtxErr
		}
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templruntime.GetBuffer(templ_7745c5c3_W)
		if !templ_7745c5c3_IsBuffer {
			defer func() {
				templ_7745c5c3_BufErr := templruntime.ReleaseBuffer(templ_7745c5c3_Buffer)
				if templ_7745c5c3_Err == nil {
					templ_7745c5c3_Err = templ_7745c5c3_BufErr
				}
			}()
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var2 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var2 == nil {
			templ_7745c5c3_Var2 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 2, "<nav class=\"bg-gray-800\"><div class=\"max-w-[140rem] mx-auto px-4 sm:px-6 lg:px-8\"><div class=\"flex h-20 items-center md:justify-between\"><!-- Logo a la izquierda --><div class=\"flex-shrink-0\"><a href=\"#\" class=\"justify-self-start text-white text-4xl font-bold 2xl:text-5xl\">GoCMS</a></div><!-- Wrap menus and dark mode button in a flex container with justify-between --><div class=\"flex items-center space-x-4 md:space-x-6 justify-between flex-1\"><!-- Menus justified to the end --><div class=\"hidden md:flex items-center justify-end space-x-4 flex-1\"><div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		for _, link := range links {
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 3, "<a class=\"text-gray-100 hover:text-gray-400 w-auto text-xl 2xl:text-3xl inline-block p-3 text-center leading-4 font-semibold transition duration-300 ease-in\" href=\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var3 templ.SafeURL
			templ_7745c5c3_Var3, templ_7745c5c3_Err = templ.JoinURLErrs(templ.URL(link.Href))
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `views/header.templ`, Line: 45, Col: 36}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var3))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 4, "\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var4 string
			templ_7745c5c3_Var4, templ_7745c5c3_Err = templ.JoinStringErrs(link.Name)
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `views/header.templ`, Line: 46, Col: 20}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var4))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 5, "</a> ")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		for name, dropdown_links := range dropdowns {
			templ_7745c5c3_Err = components.MakeDropdownButton(name, name, dropdown_links).Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 6, "</div></div><!-- Dark mode button stays on the right --></div><div class=\"md:hidden flex items-center space-x-4\"><div><button id=\"menu-toggle\" class=\"text-gray-100 dark:text-gray-100 hover:text-gray-400 focus:outline-none\"><svg class=\"w-7 h-7\" fill=\"none\" stroke=\"currentColor\" viewBox=\"0 0 24 24\"><path stroke-linecap=\"round\" stroke-linejoin=\"round\" stroke-width=\"2\" d=\"M4 6h16M4 12h16M4 18h16\"></path></svg></button></div></div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = makeDarkModeButton().Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 7, "</div></div><div id=\"mobile-menu\" class=\"hidden md:hidden\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		for _, link := range links {
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 8, "<a class=\"text-gray-100 dark:text-gray-100 hover:text-gray-400 block px-2 py-1\" href=\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var5 templ.SafeURL
			templ_7745c5c3_Var5, templ_7745c5c3_Err = templ.JoinURLErrs(templ.URL(link.Href))
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `views/header.templ`, Line: 67, Col: 111}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var5))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 9, "\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var6 string
			templ_7745c5c3_Var6, templ_7745c5c3_Err = templ.JoinStringErrs(
				link.Name)
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `views/header.templ`, Line: 69, Col: 13}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var6))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 10, "</a> ")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		for name, dropdown_links := range dropdowns {
			templ_7745c5c3_Err = components.MakeDropdownButton(name, fmt.Sprintf("%s-mobile", name), dropdown_links).Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 11, "</div></nav><hr class=\"border-t-2 border-gray-300\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return nil
	})
}

var _ = templruntime.GeneratedTemplate
