package paginator

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"math"
	"web-forum/internal/app/database/db"
	"web-forum/internal/app/models"
	"web-forum/pkg/stuff"
)

var ctx = context.Background()

func Query(preQuery models.PaginatorPreQuery) (tx pgx.Tx, rows pgx.Rows, paginatorList models.Paginator, err error) {
	var errorFunction = "paginator.Query" + preQuery.TableName

	tx, err = db.Postgres.Begin(ctx)

	if err != nil {
		return nil, nil, paginatorList, err
	}

	var topicsCount float64
	var tableName = preQuery.TableName
	var outputColumns = preQuery.OutputColumns
	var page = preQuery.Page
	var columnName = preQuery.WhereColumn
	var id = preQuery.WhereValue

	var orderStr string
	var orderDesc = preQuery.OrderReverse

	if orderDesc {
		orderStr = "desc"
	}

	if preQuery.QueryCount.PreparedValue != 0 {
		topicsCount = float64(preQuery.QueryCount.PreparedValue)
	} else {
		queryRow := tx.QueryRow(ctx, preQuery.QueryCount.Query)
		countMessagesErr := queryRow.Scan(&topicsCount)

		if countMessagesErr != nil {
			stuff.ErrLog(errorFunction, countMessagesErr)
			topicsCount = 1
		}
	}

	var pagesCount = math.Ceil(topicsCount / models.MaxPaginatorMessages)

	if float64(page) > pagesCount || float64(page) < 0 {
		page = 1
	}

	var whereStr string

	if preQuery.WhereOperator == "" {
		preQuery.WhereOperator = "="
	}

	if id != nil {
		whereStr = fmt.Sprintf("where %s %s $1", columnName, preQuery.WhereOperator)
	}

	fmtQuery := fmt.Sprintf(`select %s
	from %s
	where id in (
		select id from (
			select id, row_number() over(order by id)
			from %s
			%s
			offset %d
			limit %d
		)
		order by id %s
	)
	order by id %s;`, outputColumns, tableName, tableName, whereStr, (page-1)*models.MaxPaginatorMessages, models.MaxPaginatorMessages, orderStr, orderStr)

	if id != nil {
		rows, err = tx.Query(ctx, fmtQuery, id)
	} else {
		rows, err = tx.Query(ctx, fmtQuery)
	}

	if err != nil {
		return nil, nil, paginatorList, err
	}

	paginatorList.CurrentPage = page
	paginatorList.AllPages = int(pagesCount)

	Construct(&paginatorList)

	return tx, rows, paginatorList, nil
}

func Construct(paginatorList *models.Paginator) {
	var currentPageInt = (*paginatorList).CurrentPage
	var ourPages = (*paginatorList).AllPages
	var howMuchPagesWillBeVisible = models.HowMuchPagesWillBeVisibleInPaginator
	var dividedBy2 = float64(howMuchPagesWillBeVisible) / 2
	var floorDivided = int(math.Floor(dividedBy2))
	var ceilDivided = int(math.Ceil(dividedBy2))

	if ourPages < models.HowMuchPagesWillBeVisibleInPaginator {
		howMuchPagesWillBeVisible = ourPages
	}

	if currentPageInt > ourPages {
		currentPageInt = ourPages
	}

	currentPageInt = currentPageInt - 1 // Массив с нуля начинается.
	var limitMin, limitMax = currentPageInt - floorDivided, currentPageInt + floorDivided

	if limitMin < 0 {
		limitMin = 0
	}

	if limitMax > ourPages-1 {
		limitMax = ourPages - 1
	}

	if currentPageInt < ceilDivided {
		limitMax = howMuchPagesWillBeVisible - 1
	} else if currentPageInt >= ourPages-ceilDivided {
		limitMin = ourPages - howMuchPagesWillBeVisible
	}

	var paginatorPages = make([]int, limitMax-limitMin+1)
	var paginatorKey = 0

	for showedPage := limitMin; showedPage <= limitMax; showedPage++ {
		paginatorPages[paginatorKey] = showedPage + 1
		paginatorKey += 1
	}

	(*paginatorList).PagesArray = paginatorPages
	currentPageInt += 1

	if currentPageInt > 1 {
		(*paginatorList).Left.Activated = true
		(*paginatorList).Left.WhichPage = currentPageInt - 1
	}

	if currentPageInt < ourPages {
		(*paginatorList).Right.Activated = true
		(*paginatorList).Right.WhichPage = currentPageInt + 1
	}
}
