package airport

import (
	"context"
	"fmt"
	"log"
	"sync"
	"sync/atomic"
	"time"
)

const PlaneValue = 10
const parallelAirTrafficController = 10


// TODO:: должна быть структура долдна быть потокобезопасная
type runway struct {
	sync.Mutex
	isBusy bool
}

// useRunway испольщования взлетнопосадочной полосы
// в один моомент только один самолет может испольщовать поля данной структуры
func (r *runway) useRunway(plane *Plane, action string) {
	r.Lock()
	defer r.Unlock()
	r.isBusy = true

	plane.status = action
	fmt.Printf("Plane #%d is %s\n",plane.title,action)
	time.Sleep(time.Second * 1)

	r.isBusy = false
}

type Airport struct {
	runway

	// каналы для управления посадками и взлетами
	takeoffCh chan *Plane
	landingCh chan *Plane

	// поля которые помогут вам закрыть аэропорт
	stctx context.Context
	stop  context.CancelFunc

	// после которое говорит об завершении всех дел - программа может умирать
	done  chan struct{}
	close bool
	mx sync.Mutex
}

// isClose функция для проверки - можно ли взлетать?
func (a *Airport) isClose() bool {
	a.mx.Lock()
	defer a.mx.Unlock()
	return a.close
}

// NewAirport создание новой аэропорта и запуска его действия в отдельной горутине
func NewAirport() *Airport {
	stctx, stop := context.WithCancel(context.Background())

	a := &Airport{
		landingCh: make(chan *Plane),
		takeoffCh: make(chan *Plane),

		done: make(chan struct{}),

		stctx: stctx,
		stop:  stop,
	}

	// запускаем ассинхроную функцию для функционирования аэропорта
	//TODO:: как обработать ошибку от функции airportProcess?
	go a.airportProcess()

	return a
}

// airTrafficController создание воркера для уплавления самолетами
func (a *Airport) airTrafficController(wg *sync.WaitGroup, activePlanesCount *int64) {
	defer wg.Done()
	LOOP:
		for {
			// TODO:: проверяем есть ли неприпаркованыне самолеты
			if atomic.LoadInt64(activePlanesCount) == 0{
				break LOOP
			}
			select {
			case plane, _ := <-a.takeoffCh: // логика взлета самолета
				// проверка - можно ли совершать вылеты,
				// если нет - самолет остаётся в статусе parking и более ничего не делает
				// запрет на взлет нужно залогировать
				fmt.Printf("Plane #%d ready to takeoff\n",plane.title)

				a.mx.Lock()
				check := a.close
				a.mx.Unlock()
				if check {
					log.Printf("Plane #%d takeoff canceled, airport closed\n",plane.title)
					atomic.AddInt64(activePlanesCount,-1)
					continue
				}

				
				a.useRunway(plane, "takeoff")
				go plane.flying(a) // полетел

			case plane, _ := <-a.landingCh: // логика посадки самолета
				a.useRunway(plane, "landing")

				// обслуживаться одновременно могут только 3 самодета.
				go plane.servicing(a) // на сервисе

			default:
				time.Sleep(time.Millisecond * 10)
			}
		}
}

// airportProcess фнкция обслудивания самолетов - создание диспетчеров
func (a *Airport) airportProcess() error {
	wg := &sync.WaitGroup{}
	activePlanesCount := int64(PlaneValue)

	wg.Add(parallelAirTrafficController)
	for i := 0; i < parallelAirTrafficController; i++ {
		go a.airTrafficController(wg, &activePlanesCount)
	}

	// нужно дождаться завершения всех самолетных дел,
	// после отправить сигнал в метод Close(), что можно закрываться

	wg.Wait()
	a.done <- struct{}{}

	return nil
}

// Start запуск работы аэропорта
func (a *Airport) Start() [PlaneValue]*Plane {
	planes := [PlaneValue]*Plane{}

	for i, _ := range planes {
		planes[i] = &Plane{title: i, status: "starting"}
		go func(p *Plane) {
			a.takeoffCh <- p
		}(planes[i])
	}

	return planes
}

// Close остановку работы аэропорта
func (a *Airport) Close(seconds time.Duration) {
	time.Sleep(time.Second * seconds)

	fmt.Printf("Airport is closing...\n")
	// новые самолеты не должны вылетать, остальные должны пойти на посадку в срочном порядке
	// также остановить обслуживание если оно проходит в данный момент
	// вы обязаны дождаться пока все самолеты не закончат все свои дела
	a.stop()
	a.mx.Lock()
	a.close = true
	a.mx.Unlock()
	<- a.done
}
