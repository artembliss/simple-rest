package postgre

import (
	"fmt"
	"regexp"
	"rest-api/internal/domain"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
)

func newTestStorage(t *testing.T) (*Storage, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("unexpected error when opening stub database connection: %v", err)
	}
	t.Helper()
	sqlxDB := sqlx.NewDb(db, "postgres")
	s := &Storage{db: sqlxDB}
	return s, mock
}

func TestGetItems(t *testing.T) {
    tests := []struct {
        name          string
        mockBehavior  func(mock sqlmock.Sqlmock)
        expectedError bool
        expectedLen   int
    }{
        {
            name: "Success",
            mockBehavior: func(mock sqlmock.Sqlmock) {
                rows := sqlmock.NewRows([]string{"id", "name", "description"}).
                    AddRow(1, "Test Item", "Description")
                mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, description FROM items")).
                    WillReturnRows(rows)
            },
            expectedError: false,
            expectedLen:   1,
        },
		{
            name: "Success 2",
            mockBehavior: func(mock sqlmock.Sqlmock) {
                rows := sqlmock.NewRows([]string{"id", "name", "description"}).
                    AddRow(1, "Test Item", "Test Description 2").
                    AddRow(2, "Test Item 2", "Test Description 2")
                mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, description FROM items")).
                    WillReturnRows(rows)
            },
            expectedError: false,
            expectedLen:   2,
        },
        {
            name: "DB error",
            mockBehavior: func(mock sqlmock.Sqlmock) {
                mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, description FROM items")).
                    WillReturnError(fmt.Errorf("db failure"))
            },
            expectedError: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            s, mock := newTestStorage(t)      			
			t.Cleanup(func(){s.db.Close()})


            tt.mockBehavior(mock)

            items, err := s.GetItems()
            if tt.expectedError && err == nil {
                t.Errorf("expected error, got nil")
				return
            }
            if !tt.expectedError && err != nil {
                t.Errorf("unexpected error: %v", err)
				return
            }
            if len(items) != tt.expectedLen {
                t.Errorf("expected %d items, got %d", tt.expectedLen, len(items))
				return
            }
            if err := mock.ExpectationsWereMet(); err != nil {
                t.Errorf("unmet expectations: %v", err)
				return
            }
        })
    }
}


func TestGetItem(t *testing.T){
	tests := []struct{
		name          string
		expectedItem domain.Item
		inputID int
		mockBehavior  func(mock sqlmock.Sqlmock)
		expectedError bool
	}{
		{
			name: "Success",
			expectedItem: domain.Item{ID: 1, Name: "Test name", Description: "Test description"},
			inputID: 1,
			mockBehavior: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "name", "description"}).
				AddRow(1, "Test name", "Test description")
				mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, description FROM items WHERE id = $1")).
				WithArgs(1).
				WillReturnRows(rows)
			},
			expectedError: false,
		},
		{
			name: "Non-existed id",
			inputID: 42,
			expectedItem: domain.Item{},
			mockBehavior: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "name", "description"})
				mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, description from items WHERE id = $1")).
				WithArgs(42).WillReturnRows(rows)
			},
			expectedError: true,
		},
		{
			name:         "DB error",
			expectedItem: domain.Item{},
			mockBehavior: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, description FROM items WHERE id = $1")).
					WithArgs(1).
					WillReturnError(fmt.Errorf("db failure"))
			},
			expectedError: true,
		},
		{
			name:         "Scan error",
			expectedItem: domain.Item{},
			inputID: 1,
			mockBehavior: func(mock sqlmock.Sqlmock) {
				// Симулируем ошибку сканирования: передаем неверный тип для поля id.
				rows := sqlmock.NewRows([]string{"id", "name", "description"}).
					AddRow("invalid", "Test Item", "Description")
				mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, description FROM items WHERE id = $1")).
					WithArgs(1).
					WillReturnRows(rows)
			},
			expectedError: true,
		},
		{
			name: "Zero id",
			inputID: 0,
			mockBehavior: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "name", "description"}).
				AddRow(1, "Test name", "Test description")
				mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, description FROM items WHERE id = $1")).
				WithArgs(0).
				WillReturnError(fmt.Errorf("failed to get item: sql: no rows in result set"))
				_ = rows
			},
			expectedItem: domain.Item{},
			expectedError: true,
		},
	}
	for _, tt := range tests{
		t.Run(tt.name, func(t *testing.T) {	
			s, mock := newTestStorage(t)
			t.Cleanup(func(){s.db.Close()})


			tt.mockBehavior(mock)

			item, err := s.GetItem(tt.inputID)
			if tt.expectedError {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
				return
			}
            if !tt.expectedError && err != nil {
                t.Errorf("unexpected error: %v", err)
				return
			}
			if item != tt.expectedItem {
				t.Errorf("expected item %+v, got %+v", tt.expectedItem, item)
				return
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unmet expectations: %v", err)
				return
			}
		})
	}
}

func TestCreateItem(t *testing.T){
	tests := []struct{
		name string
		expectedItem domain.Item
		mockBehavior func(mock sqlmock.Sqlmock)
		expectedError bool
	}{
		{
			name: "Success",
			expectedItem: domain.Item{ID: 1, Name: "Test Name", Description: "Test Description"},
			mockBehavior: func(mock sqlmock.Sqlmock){
				rows := sqlmock.NewRows([]string{"id", "name", "description"}).
				AddRow(1, "Test Name", "Test Description")
				mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO items(name, description) VALUES($1, $2)")).
				WithArgs("Test Name", "Test Description").
				WillReturnRows(rows)
			},
			expectedError: false,
		},
		{
			name: "Scan Error",
			expectedItem: domain.Item{ID: 1, Name: "Test Name", Description: "Test Description"},
			mockBehavior: func(mock sqlmock.Sqlmock){
				rows := sqlmock.NewRows([]string{"id", "name", "description"}).
				AddRow("invalid_id", "Test Name", "Test Description")
				mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO items(name, description) VALUES($1, $2)")).
				WithArgs("Test Name", "Test Description").
				WillReturnRows(rows)
			},
			expectedError: true,
		},
		{
			name: "Storage Error",
			expectedItem: domain.Item{ID: 1, Name: "Test Name", Description: "Test Description"},
			mockBehavior: func(mock sqlmock.Sqlmock){
				mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO items(name, description) VALUES($1, $2)")).
				WithArgs("Test Name", "Test Description").
				WillReturnError(fmt.Errorf("failed to create item"))
			},
			expectedError: true,
		},
		{
			name: "Invalid input",
			expectedItem: domain.Item{ID: 1, Name: "Test Name", Description: "Test Description"},
			mockBehavior: func(mock sqlmock.Sqlmock){
				rows := sqlmock.NewRows([]string{"id", "name", "description"}).
				AddRow("1", "Test Name", "Test Description")
				mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO items(name, description) VALUES($1, $2)")).
				WithArgs("", "").
				WillReturnError(fmt.Errorf("failed to create item"))
				_ = rows
			},
			expectedError: true,
		},
		
	}
	for _, tt := range tests{
		t.Run(tt.name, func(t *testing.T) {
			s, mock := newTestStorage(t)
			t.Cleanup(func(){s.db.Close()})


			tt.mockBehavior(mock)
			
			item, err := s.CreateItem(tt.expectedItem)
			if tt.expectedError {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
				return
			}
			if !tt.expectedError && err != nil{
				t.Errorf("unexpected error: %v", err)
				return
			} 
			if item != tt.expectedItem{
				t.Errorf("expected item: %v, got: %v", tt.expectedItem, item)
				return
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unmet expectations: %v", err)
				return
			}
		})
	}
}
//Доделать mockBehavior и добавить новых тестов, разобраться в логике работы
func TestUpdateItem(t *testing.T){
	tests := []struct{
		name string
		expectedItem domain.Item
		mockBehavior func(mock sqlmock.Sqlmock)
		expectedError bool
	}{
		{
			name: "Success",
			expectedItem: domain.Item{ID: 1, Name: "Updated Test Name", Description: "Updated Test Description"},
			mockBehavior: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "name", "description"}).
				AddRow(1, "Updated Test Name", "Updated Test Description")
				mock.ExpectQuery(regexp.QuoteMeta("UPDATE items SET name = $1, description = $2 WHERE id = $3")).
				WithArgs("Updated Test Name", "Updated Test Description", 1).
				WillReturnRows(rows)
			},
			expectedError: false,
		},
		{
			name: "Scan error",
			expectedItem: domain.Item{ID: 1, Name: "Updated Test Name", Description: "Updated Test Description"},
			mockBehavior: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "name", "description"}).
				AddRow("invalid_id", "Updated Test Name", "Updated Test Description")
				mock.ExpectQuery(regexp.QuoteMeta("UPDATE items SET name = $1, description = $2 WHERE id = $3")).
				WithArgs("Updated Test Name", "Updated Test Description", 1).
				WillReturnRows(rows)

			},
			expectedError: true,
		},
		{
			name: "Storage error",
			expectedItem: domain.Item{ID: 1, Name: "Updated Test Name", Description: "Updated Test Description"},
			mockBehavior: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("UPDATE items SET name = $1, description = $2 WHERE id = $3")).
				WithArgs("Updated Test Name", "Updated Test Description", 1).
				WillReturnError(fmt.Errorf("failed to update item"))

			},
			expectedError: true,
		},
		{
			name: "Input error",
			expectedItem: domain.Item{ID: 1, Name: "Updated Test Name", Description: "Updated Test Description"},
			mockBehavior: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("UPDATE items SET name = $1, description = $2 WHERE id = $3")).
				WithArgs("", "", 1).
				WillReturnError(fmt.Errorf("failed to update item"))

			},
			expectedError: true,
		},
	}
	for _, tt := range tests{
		t.Run(tt.name, func(t *testing.T) {	
			s, mock := newTestStorage(t)
			t.Cleanup(func(){s.db.Close()})


			tt.mockBehavior(mock)

			item, err := s.UpdateItem(tt.expectedItem.ID, tt.expectedItem)
			if tt.expectedError {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
				return
			}
            if !tt.expectedError && err != nil {
                t.Errorf("unexpected error: %v", err)
				return
			}
			if item != tt.expectedItem {
				t.Errorf("expected item %+v, got %+v", tt.expectedItem, item)
				return
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unmet expectations: %v", err)
				return
			}
		})
	}
}

func TestDeleteItem(t *testing.T){
	tests := []struct{
		name string
		expectedItem domain.Item
		mockBehavior func(mock sqlmock.Sqlmock)
		expectedError bool
	}{
		{
			name: "Success",
			expectedItem: domain.Item{ID: 1, Name: "Delete Test Name", Description: "Delete Test Description"},
			mockBehavior: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "name", "description"}).
				AddRow(1, "Delete Test Name", "Delete Test Description")
				mock.ExpectQuery(regexp.QuoteMeta("DELETE FROM items WHERE id = $1 RETURNING id, name, description")).
				WithArgs(1).
				WillReturnRows(rows)
				mock.ExpectExec(regexp.QuoteMeta("DELETE id, name, description FROM items WHERE id = $1")).
				WithArgs(1).
				WillReturnResult(sqlmock.NewResult(0, 1))
			},
			expectedError: false,
		},
	}
	for _, tt := range tests{
		t.Run(tt.name,func(t *testing.T) {
			s, mock := newTestStorage(t)
			t.Cleanup(func(){s.db.Close()})

			tt.mockBehavior(mock)

			item, err := s.DeleteItem(tt.expectedItem.ID)

			if tt.expectedError{
				if err == nil{
					t.Errorf("expected error, got nil")
				}
				return
			}
			if !tt.expectedError && err != nil{
                t.Errorf("unexpected error: %v", err)
				return
			}
			if item != tt.expectedItem{
				t.Errorf("expected item %+v, got %+v", tt.expectedItem, item)
				return
			}
			if err := mock.ExpectationsWereMet(); err != nil{
				t.Errorf("unmet expectations: %v", err)
				return
			}
		})
	}
}