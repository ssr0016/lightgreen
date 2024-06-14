package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Graceful shutdown... Intercepting Shutdown Signals ... Executing the shutdown
func (app *application) server() error {
	// Declare a HTTP server using the same settings as in main()
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.config.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	// create a shutdownError channel
	shutdownError := make(chan error)

	// Start a background goroutine
	go func() {
		// Intercept the signals
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit

		// Update the log entry to say "shutting down server" instead of "caught signal"
		app.logger.PrintInfo("shutting down server", map[string]string{
			"signal": s.String(),
		})

		// Create a context with a 5-second timeout
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Call Shutdown() on the server as before, sending on the shutdownError channel only if it returns an error.
		err := srv.Shutdown(ctx)
		if err != nil {
			shutdownError <- err
		}

		// Log a message indicating that we're waiting for background goroutines to finish their tasks.
		app.logger.PrintInfo("completing background tasks", map[string]string{
			"addr": srv.Addr,
		})

		// Wait for background goroutines to finish using Wait(), then return nil
		// on the shutdownError channel to indicate successful shutdown.
		app.wg.Wait()
		shutdownError <- nil
	}()

	app.logger.PrintInfo("starting server", map[string]string{
		"addr": srv.Addr,
		"env":  app.config.env,
	})

	// Shutdown() causes ListenAndServe() to return http.ErrServerClosed,
	// indicating graceful shutdown has started. Return only if it's NOT http.ErrServerClosed.
	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	// Wait for the return value from Shutdown() on shutdownError.
	// Return the error if there was a problem with graceful shutdown.
	err = <-shutdownError
	if err != nil {
		return err
	}

	// Graceful shutdown completed successfully, log "stopped server".
	app.logger.PrintInfo("stopped server", map[string]string{
		"addr": srv.Addr,
	})

	return nil
}
