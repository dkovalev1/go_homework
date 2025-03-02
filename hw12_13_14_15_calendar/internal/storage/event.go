package storage

import "time"

/*
Событие - основная сущность, содержит в себе поля:
* ID - уникальный идентификатор события (можно воспользоваться UUID);
* Заголовок - короткий текст;
* Дата и время события;
* Длительность события (или дата и время окончания);
* Описание события - длинный текст, опционально;
* ID пользователя, владельца события;
* За сколько времени высылать уведомление, опционально.
*/
type Event struct {
	ID          string        `db:"id"`
	Title       string        `db:"title"`
	StartTime   time.Time     `db:"starttime"`
	Duration    time.Duration `db:"duration"`
	Description string        `db:"description"`
	UserID      int64         `db:"userid"`
	NotifyTime  time.Duration `db:"notifytime"`
}
