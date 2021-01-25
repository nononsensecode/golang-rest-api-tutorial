package model

import (
	"database/sql"
	"nononsensecode/rest-api-tutorial/articleerror"
)

// DB Database object
var DB *sql.DB

// Article type exposes our article object
type Article struct {
	ID      int64 `json:"Id"`
	Title   string `json:"Title"`
	Desc    string `json:"desc"`
	Content string `json:"content"`
}

//Articles slice of articles
var Articles []Article

// CreateTable creates the article table
func CreateTable() error {
	statement, err := DB.Prepare("CREATE TABLE IF NOT EXISTS article(id INTEGER PRIMARY KEY AUTOINCREMENT, title TEXT, desc TEXT, content TEXT)")
	if err != nil {
		return err		
	}

	defer statement.Close()

	statement.Exec()

	return nil
}

// CreateArticle creates a new row in the article table
func CreateArticle(article Article) (int64, error) {
	transaction, err := DB.Begin()
	if err != nil {
		return 0, err
	}
	defer transaction.Rollback()

	statement, err := DB.Prepare("INSERT INTO article (title, content, desc) VALUES (?, ?, ?)")
	if err != nil {
		return 0, err
	}
	defer statement.Close()

	result, err := statement.Exec(article.Title, article.Content, article.Desc)
	if err != nil {
		return 0, err
	}
	
	err = transaction.Commit()
	if err != nil {
		return 0, err
	}

	article.ID, err = result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return article.ID, nil
}

// FindArticleByID return an article by a specified id or nil
func FindArticleByID(id int64) (Article, error) {
	var article Article
	err := DB.QueryRow("SELECT * FROM article WHERE id = ?", id).
	Scan(&article.ID, &article.Title, &article.Content, &article.Desc)
	
	switch {
	case err == sql.ErrNoRows:
		return Article{}, articleerror.ArticleDoesNotExistError(id)
	case err != nil:
		return Article{}, articleerror.UnknownError(err)
	default:
		return article, nil
	}
}

// FindAllArticles returns all articles in the table
func FindAllArticles() ([]Article, error) {
	var articles []Article

	rows, err := DB.Query("SELECT * FROM article")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var article Article

		err = rows.Scan(&article.ID, &article.Title, &article.Content, &article.Desc)
		if err != nil {
			return nil, err
		}

		articles = append(articles, article)
	}

	return articles, nil
}

// DeleteArticleByID deletes an article according to the ID supplied
func DeleteArticleByID(id int64) error {
	transaction, err := DB.Begin()
	if err != nil {
		return err
	}
	defer transaction.Rollback()

	statement, err := DB.Prepare("DELETE FROM article WHERE id = ?")
	if err != nil {
		return err
	}
	defer statement.Close()

	_, err = statement.Exec()
	if err != nil {
		return err
	}

	err = transaction.Commit()
	if err != nil {
		return err
	}

	return nil
}

// UpdateArticle updates article by ID
func UpdateArticle(article Article) (Article, error) {
	transaction, err := DB.Begin()
	if err != nil {
		return Article{}, err
	}
	defer transaction.Rollback()

	statement, err := DB.Prepare("UPDATE article SET title = ?, content = ?, desc = ? WHERE id = ?")
	if err != nil {
		return Article{}, err
	}
	defer statement.Close()

	result, err := statement.Exec(article.Title, article.Content, article.Desc, article.ID)
	if err != nil {
		return Article{}, err
	}

	err = transaction.Commit()
	if err != nil {
		return Article{}, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return Article{}, err
	}

	if rowsAffected > 0 {
		return article, nil
	} else {
		return Article{}, err
	}
}