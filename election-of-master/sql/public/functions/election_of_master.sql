CREATE OR REPLACE FUNCTION election_of_master
(
	_table_name varchar(255),
	_service_instance_id varchar(255),
	_expiration_period_sec integer
)
RETURNS RECORD
LANGUAGE plpgsql
AS

$$
DECLARE
	_now timestamp with time zone;
	_lat_limit timestamp with time zone;
    _lat_new timestamp with time zone;
    _current_service_instance_id varchar(255);
    _current_service_instance_id_new varchar(255);
    _current_last_update_time timestamp with time zone;
    _last_update_time_feedback timestamp with time zone;
    _cmd text;
    _result RECORD;
BEGIN
    --RAISE NOTICE '[INPUT] SIID:% EPS:% NOW:%', _service_instance_id, _expiration_period_sec, now(); --// Debug //--
	_now = now();

	--// Get Data from the current Master Record //--
	_cmd = 'SELECT "ServiceInstanceId", "LastUpdateTime" FROM "' || _table_name || '";';
    EXECUTE _cmd INTO _current_service_instance_id, _current_last_update_time;
    --RAISE NOTICE '[MASTER RECORD] ID:% _LAT:%', _current_service_instance_id, _current_last_update_time; --// Debug //--

    --// Get the current Master Record and check it for Consistency //--
	IF (_current_service_instance_id IS NULL) OR
	   (_current_last_update_time IS NULL) THEN
	BEGIN
        RAISE NOTICE 'Master Record is damaged or does not exist.';

        --// Reset the Table //--
        _cmd = 'DELETE FROM "' || _table_name || '" WHERE TRUE;';
        EXECUTE _cmd;

        --// Take the Lease //--
        _lat_new = _now;
        _current_service_instance_id_new = _service_instance_id;
        _cmd =  'INSERT INTO "' || _table_name || '" ("ServiceInstanceId", "LastUpdateTime") ' ||
                'VALUES (''' || _current_service_instance_id_new || ''', ''' || _lat_new || '''::timestamptz) ' ||
                'RETURNING "LastUpdateTime";';
        EXECUTE _cmd INTO _last_update_time_feedback;
        
        SELECT TRUE, _last_update_time_feedback INTO _result;
        RETURN _result;
    END;
	END IF;

	_lat_limit = _current_last_update_time + (interval '1 second' * _expiration_period_sec);
    IF (_now <= _lat_limit) THEN
    BEGIN

        IF (_current_service_instance_id = _service_instance_id) THEN
        BEGIN
            --RAISE NOTICE 'Renewing the Lease.'; --// Debug //--

            --// Reset the Table //--
            _cmd = 'DELETE FROM "' || _table_name || '" WHERE TRUE;';
            EXECUTE _cmd;

            --// Renew the Lease //--
            _lat_new = _now;
            _current_service_instance_id_new = _service_instance_id;
            _cmd =  'INSERT INTO "' || _table_name || '" ("ServiceInstanceId", "LastUpdateTime") ' ||
                'VALUES (''' || _current_service_instance_id_new || ''', ''' || _lat_new || '''::timestamptz) ' ||
                'RETURNING "LastUpdateTime";';
			EXECUTE _cmd INTO _last_update_time_feedback;
        
			SELECT TRUE, _last_update_time_feedback INTO _result;
			RETURN _result;
        END;
        ELSE
        BEGIN
            --RAISE NOTICE 'Lease is busy.'; --// Debug //--

            SELECT FALSE, _current_last_update_time INTO _result;
			RETURN _result;
        END;
        END IF;

    END;
    ELSE
    BEGIN
        --RAISE NOTICE 'Taking the Lease.'; --// Debug //--

         --// Reset the Table //--
        _cmd = 'DELETE FROM "' || _table_name || '" WHERE TRUE;';
        EXECUTE _cmd;

        --// Take the Lease //--
        _lat_new = _now;
        _current_service_instance_id_new = _service_instance_id;
        _cmd =  'INSERT INTO "' || _table_name || '" ("ServiceInstanceId", "LastUpdateTime") ' ||
                'VALUES (''' || _current_service_instance_id_new || ''', ''' || _lat_new || '''::timestamptz) ' ||
                'RETURNING "LastUpdateTime";';
        EXECUTE _cmd INTO _last_update_time_feedback;
        
        SELECT TRUE, _last_update_time_feedback INTO _result;
        RETURN _result;
    END;
    END IF;

    --// For any Case //--
	SELECT FALSE, _current_last_update_time INTO _result;
	RETURN _result;
END
$$;
