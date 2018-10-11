package main

import (
	"strconv"

	"github.com/labstack/echo"
	"github.com/tomnomnom/linkheader"
)

// ParsePagination 获得页码，每页条数，Echo中间件。
func ParsePagination(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var err error
		var page, pageSize int
		// 获得页码
		if c.QueryParam("page") == "" {
			page = 1
		} else {
			if page, err = strconv.Atoi(c.QueryParam("page")); err != nil {
				return newHTTPError(400, "InvalidPage", "请在URL中提供合法的页码")
			}
		}
		// 获得每页条数
		if c.QueryParam("per_page") == "" {
			pageSize = config.APP.PageSize
		} else {
			if pageSize, err = strconv.Atoi(c.QueryParam("per_page")); err != nil {
				return newHTTPError(400, "InvalidPage", "请在URL中提供合法的每页条数")
			}
		}
		// 设置查询数据时的 offset 和 limit
		c.Set("page", page)
		c.Set("offset", (page-1)*pageSize)
		c.Set("limit", pageSize)
		return next(c)
	}
}

// setPaginationHeader 设置分页相关 resp header
// 如果要显示页码，还需要返回 X-Total-Count 和 Link 的 last 信息，可以多传入一个记录总数参数进行处理。
// 移动应用一般不用知道总条数，传统的web分页器有时会需要。
func setPaginationHeader(c echo.Context, isLast bool) {
	page := c.Get("page").(int)
	pageSize := c.Get("limit").(int)
	c.Response().Header().Set("X-Page-Num", strconv.Itoa(page))
	c.Response().Header().Set("X-Page-Size", strconv.Itoa(pageSize))
	link := linkheader.Links{
		{URL: config.APP.BaseURL + "?page=" + strconv.Itoa(page) + "&per_page=" + strconv.Itoa(pageSize), Rel: "self"},
	}
	if !isLast {
		link = append(link, linkheader.Link{URL: config.APP.BaseURL + "?page=" + strconv.Itoa(page+1) + "&per_page=" + strconv.Itoa(pageSize), Rel: "next"})
	}
	c.Response().Header().Set("Link", link.String())
	return
}
