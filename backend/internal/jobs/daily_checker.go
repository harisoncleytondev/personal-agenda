package jobs

import (
	"context"
	"log"
	"time"

	"github.com/harisoncleytondev/personal-agenda/internal/service"
	"github.com/robfig/cron/v3"
)

func Start(s *service.AppointmentService) {
	loc, _ := time.LoadLocation("America/Sao_Paulo")

	c := cron.New(cron.WithLocation(loc))

	_, err := c.AddFunc("0 5 * * *", func() {
		log.Println("[CRON] Iniciando rotina diária de agendamentos...")
		
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cancel()

		err := s.ProcessDailyRoutines(ctx)
		if err != nil {
			log.Printf("[CRON] Erro ao processar rotina diária: %v\n", err)
			return
		}
		
		log.Println("[CRON] Rotina diária finalizada com sucesso.")
	})

	if err != nil {
		log.Fatal(err)
	}

	c.Start()
}