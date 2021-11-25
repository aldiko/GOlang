package airport

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

const (
	IsFly="fly"
	IsService="on service"
	IsParking="parking"
)

type Plane struct {
	title  int
	status string
}

// flying функция для полета самолета, в конце она отправляет самолет на посадку
func (p *Plane) flying(a *Airport) {
	p.status = IsFly

	r := rand.Intn(10)
	if r < 3 {
		r = 3
	}

	ctx, cancel :=context.WithTimeout(a.stctx,time.Second*time.Duration(r))
	defer cancel()
	fmt.Printf("Plane #%d is flying\n",p.title)
	//логика полета.
	// Полет либо должен закончится по таймауту, либо если аэропорт скажет садить - мы закрываемся
	LOOP:
		for  {
			select {
			case <-ctx.Done():
				fmt.Printf("Plane #%d is landing\n",p.title)
				break LOOP
			default:
				// для экономии ресурсов, чтобы не проверял каждый момент
				time.Sleep(time.Millisecond*10)
			}
		}

	//самолет нужно отправить на посадку
	a.landingCh <- p
}

// servicing функция обслуживания самолета, в конце она отправляет самолет обратно на взлет
func (p *Plane) servicing(a *Airport) {
	p.status = IsService

	r := rand.Intn(3)
	if r < 1 {
		r = 1
	}

	ctx , cancel := context.WithTimeout(context.Background(),time.Second*time.Duration(r))
	defer cancel()

	fmt.Printf("Plane #%d on service\n",p.title)
	// логика обслуживания самолета.
	// Обслуивание либо должено закончится по таймауту, либо если аэропорт скажет заканчивай - мы закрываемся
	LOOP:
		for {
			select {
			case <-ctx.Done():
				fmt.Printf("Plane #%d is finished service\n",p.title)
				break LOOP
			case <-a.stctx.Done():
				fmt.Printf("Plane #%d cancelled service\n",p.title)
				break LOOP
			default:
				// для экономии ресурсов
				time.Sleep(time.Millisecond*10)
			}
		}

	p.status = IsParking
	// самолет нужно отправить на попытку взлета
	a.takeoffCh <- p
}
