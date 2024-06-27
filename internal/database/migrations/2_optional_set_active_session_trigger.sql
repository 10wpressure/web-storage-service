CREATE OR REPLACE FUNCTION set_active_session()
    RETURNS TRIGGER AS $$
BEGIN
    -- Устанавливаем все предыдущие сессии пользователя как неактивные
    UPDATE sessions
    SET active = FALSE
    WHERE uid = NEW.uid;

    -- Устанавливаем текущую сессию как активную
    NEW.active = TRUE;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER  activate_latest_session
    BEFORE INSERT ON sessions
    FOR EACH ROW
EXECUTE FUNCTION set_active_session();