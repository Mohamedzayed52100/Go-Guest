package database

func SetupReservationListener(tx Executable) error {
	createFunctionSQL := `
	CREATE OR REPLACE FUNCTION notify_reservations_change()
	RETURNS TRIGGER AS $$
	DECLARE
		  payload JSON;
		  notification JSON;
		  notificationChannel TEXT := 'reservations_change';
	BEGIN	
		IF TG_OP = 'INSERT' THEN
			payload := json_build_object(
				'operation', TG_OP,
				'row', NEW
			);
		ELSIF TG_OP = 'UPDATE' THEN
			 IF NEW.deleted_at IS NULL THEN
				  payload := json_build_object(
					   'operation', TG_OP,
					   'row', NEW
				);
			ELSE
				  RETURN NULL;
			END IF;
		ELSE
			RETURN NULL;
		END IF;
	
		notification := json_build_object(
			'channel', notificationChannel,
			'payload', payload
		);
	
		PERFORM pg_notify(notificationChannel, notification::text);
	
		RETURN NULL;
	END;
	$$ LANGUAGE plpgsql;
	`

	if _, err := tx.Exec(createFunctionSQL); err != nil {
		return err
	}

	// Drop trigger SQL
	dropTriggerSQL := `
	DROP TRIGGER IF EXISTS reservations_trigger ON reservations;
	`

	if _, err := tx.Exec(dropTriggerSQL); err != nil {
		return err
	}

	// Create trigger SQL
	createTriggerSQL := `
	CREATE TRIGGER reservations_trigger
	AFTER INSERT OR UPDATE ON reservations
	FOR EACH ROW EXECUTE PROCEDURE notify_reservations_change();
	`

	if _, err := tx.Exec(createTriggerSQL); err != nil {
		return err
	}

	return nil
}

func SetupReservationsWaitlistListener(tx Executable) error {
	createFunctionSQL := `
	CREATE OR REPLACE FUNCTION notify_reservation_waitlists_change()
	RETURNS TRIGGER AS $$
	DECLARE
		  payload JSON;
		  notification JSON;
		  notificationChannel TEXT := 'reservation_waitlists_change';
	BEGIN	
		IF TG_OP = 'INSERT' THEN
			payload := json_build_object(
				'operation', TG_OP,
				'row', NEW
			);
		ELSIF TG_OP = 'UPDATE' THEN
			payload := json_build_object(
				'operation', TG_OP,
				'row', NEW
			);
		ELSE
			RETURN NULL;
		END IF;
	
		notification := json_build_object(
			'channel', notificationChannel,
			'payload', payload
		);
	
		PERFORM pg_notify(notificationChannel, notification::text);
	
		RETURN NULL;
	END;
	$$ LANGUAGE plpgsql;
	`

	if _, err := tx.Exec(createFunctionSQL); err != nil {
		return err
	}

	// Drop trigger SQL
	dropTriggerSQL := `
	DROP TRIGGER IF EXISTS reservation_waitlists_trigger ON reservation_waitlists;
	`

	if _, err := tx.Exec(dropTriggerSQL); err != nil {
		return err
	}

	// Create trigger SQL
	createTriggerSQL := `
	CREATE TRIGGER reservation_waitlists_trigger
	AFTER INSERT OR UPDATE ON reservation_waitlists
	FOR EACH ROW EXECUTE PROCEDURE notify_reservation_waitlists_change();
	`

	if _, err := tx.Exec(createTriggerSQL); err != nil {
		return err
	}

	return nil
}
