// 公众号文章列表
package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo"
)

// Note 纸条
type Note struct {
	ID int `json:"id" gorm:"primary_key"`
	// 所属用户
	UserID int `json:"user_id" gorm:"index:idx_user_update"`
	// 标题
	Title string `json:"title"`
	// 内容
	Content string `json:"content" gorm:"size:2000"`
	// 是否公开
	IsPublic bool `json:"is_public"`
	// 创建时间
	CreatedTime time.Time `json:"created_time,omitempty"`
	// 最后更新时间
	UpdatedTime time.Time  `json:"updated_time,omitempty" gorm:"index:idx_user_update"`
	DeletedAt   *time.Time `json:"-"`
}

// NoteUpdate 更新请求结构体，用指针可以判断是否有请求这个字段
type NoteUpdate struct {
	// 标题
	Title *string `json:"title"`
	// 内容
	Content *string `json:"content"`
	// 是否公开
	IsPublic *bool `json:"is_public"`
}

func findNoteByID(id int) (*Note, error) {
	var n = new(Note)
	if err := db.First(n, id).Error; err != nil {
		return nil, err
	}
	return n, nil
}

// createNote 新建笔记
// @Tags 笔记
// @Summary 新建笔记
// @Description 新建一条笔记
// @Accept  json
// @Produce  json
// @Param data body main.Note true "笔记内容"
// @Success 201 {object} main.Note
// @Failure 400 {object} main.httpError
// @Failure 401 {object} main.httpError
// @Failure 500 {object} main.httpError
// @Security ApiKeyAuth
// @Router /notes [post]
func createNote(c echo.Context) error {
	var a = new(Note)
	if err := c.Bind(a); err != nil {
		return err
	}
	// 校验
	if a.Title == "" {
		return newHTTPError(400, "BadRequest", "Empty title")
	}
	if a.Content == "" {
		return newHTTPError(400, "BadRequest", "Empty content")
	}
	// 用户信息
	userID, err := parseUser(c)
	if err != nil {
		return err
	}
	a.UserID = userID
	// 保存
	if err := db.Create(a).Error; err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, a)
}

// updateNote 更新笔记
// @Tags 笔记
// @Summary 更新笔记
// @Description 更新指定id的笔记
// @Accept  json
// @Produce  json
// @Param data body main.NoteUpdate true "更新内容"
// @Success 200 {object} main.Note
// @Failure 400 {object} main.httpError
// @Failure 401 {object} main.httpError
// @Failure 403 {object} main.httpError
// @Failure 404 {object} main.httpError
// @Failure 500 {object} main.httpError
// @Security ApiKeyAuth
// @Router /notes/{id} [put]
func updateNote(c echo.Context) error {
	// 获取URL中的ID
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return newHTTPError(400, "InvalidID", "请在URL中提供合法的ID")
	}
	var n = new(NoteUpdate)
	if err := c.Bind(n); err != nil {
		return err
	}
	old, err := findNoteByID(id)
	if err != nil {
		return err
	}
	// 用户权限
	userID, err := parseUser(c)
	if err != nil {
		return err
	}
	if userID != old.UserID {
		return ErrForbidden
	}
	// 利用指针检查是否有请求这个字段
	if n.Title != nil {
		if *n.Title == "" {
			return newHTTPError(400, "BadRequest", "Empty title")
		}
		old.Title = *n.Title
	}
	if n.Content != nil {
		if *n.Content == "" {
			return newHTTPError(400, "BadRequest", "Empty content")
		}
		old.Content = *n.Content
	}
	if n.IsPublic != nil {
		old.IsPublic = *n.IsPublic
	}

	if err := db.Save(old).Error; err != nil {
		return err
	}

	return c.JSON(http.StatusOK, old)
}

// deleteNote 删除笔记
// @Tags 笔记
// @Summary 删除笔记
// @Description 删除指定id的笔记
// @Accept  json
// @Produce  json
// @Param id path int true "笔记编号"
// @Success 204
// @Failure 400 {object} main.httpError
// @Failure 401 {object} main.httpError
// @Failure 403 {object} main.httpError
// @Failure 404 {object} main.httpError
// @Failure 500 {object} main.httpError
// @Security ApiKeyAuth
// @Router /notes/{id} [delete]
func deleteNote(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return newHTTPError(400, "InvalidID", "请在URL中提供合法的ID")
	}
	// 查询对象
	n, err := findNoteByID(id)
	if err != nil {
		return err
	}
	// 用户权限
	userID, err := parseUser(c)
	if err != nil {
		return err
	}
	if userID != n.UserID {
		return ErrForbidden
	}
	// 删除数据库对象
	if err := db.Delete(&Note{ID: id}).Error; err != nil {
		return err
	}
	return c.NoContent(http.StatusNoContent)
}

// getNote 获取笔记
// @Tags 笔记
// @Summary 获取笔记
// @Description 获取指定id的笔记
// @Accept  json
// @Produce  json
// @Param id path int true "笔记编号"
// @Success 200 {object} main.Note
// @Failure 400 {object} main.httpError
// @Failure 401 {object} main.httpError
// @Failure 403 {object} main.httpError
// @Failure 404 {object} main.httpError
// @Failure 500 {object} main.httpError
// @Security ApiKeyAuth
// @Router /notes/{id} [get]
func getNote(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return newHTTPError(400, "InvalidID", "请在URL中提供合法的ID")
	}
	n, err := findNoteByID(id)
	if err != nil {
		return err
	}
	// 用户权限
	userID, err := parseUser(c)
	if err != nil {
		return err
	}
	if userID != n.UserID && !n.IsPublic {
		return ErrForbidden
	}
	return c.JSON(http.StatusOK, n)
}

// getPublicNote 获取公开笔记
// @Tags 笔记
// @Summary 获取公开笔记
// @Description 获取指定id的公开笔记
// @Accept  json
// @Produce  json
// @Param id path int true "笔记编号"
// @Success 200 {object} main.Note
// @Failure 400 {object} main.httpError
// @Failure 404 {object} main.httpError
// @Failure 500 {object} main.httpError
// @Router /public/notes/{id} [get]
func getPublicNote(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return newHTTPError(400, "InvalidID", "请在URL中提供合法的ID")
	}
	n, err := findNoteByID(id)
	if err != nil {
		return err
	}
	if !n.IsPublic {
		return ErrNotFound
	}
	return c.JSON(http.StatusOK, n)
}

// getNotes 获取用户笔记列表
// @Tags 笔记
// @Summary 获取用户笔记列表
// @Description 获取用户的全部笔记，有分页，默认一页10条。
// @Accept  json
// @Produce  json
// @Param page query int false "页码"
// @Param per_page query int false "每页几条"
// @Success 200 {array} main.Note
// @Failure 400 {object} main.httpError
// @Failure 401 {object} main.httpError
// @Failure 500 {object} main.httpError
// @Security ApiKeyAuth
// @Router /notes [get]
func getNotes(c echo.Context) error {
	// 提前make可以让查询没有结果的时候返回空列表
	var ns = make([]*Note, 0)
	// 用户信息
	userID, err := parseUser(c)
	if err != nil {
		return err
	}
	// 分页信息
	limit := c.Get("limit").(int)
	offset := c.Get("offset").(int)
	if err := db.Where("user_id = ?", userID).Order("updated_at desc").Offset(offset).Limit(limit).Find(&ns).Error; err != nil {
		return err
	}
	return c.JSON(http.StatusOK, ns)
}

// getPublicNotes 获取公开笔记列表
// @Tags 笔记
// @Summary 获取公开笔记列表
// @Description 获取公开的全部笔记，有分页，默认一页10条。
// @Accept  json
// @Produce  json
// @Param page query int false "页码"
// @Param per_page query int false "每页几条"
// @Success 200 {array} main.Note
// @Failure 400 {object} main.httpError
// @Failure 500 {object} main.httpError
// @Router /public/notes [get]
func getPublicNotes(c echo.Context) error {
	// 提前make可以让查询没有结果的时候返回空列表
	var ns = make([]*Note, 0)
	// 分页信息
	limit := c.Get("limit").(int)
	offset := c.Get("offset").(int)
	if err := db.Where("is_public = true").Order("updated_at desc").Offset(offset).Limit(limit).Find(&ns).Error; err != nil {
		return err
	}
	return c.JSON(http.StatusOK, ns)
}
