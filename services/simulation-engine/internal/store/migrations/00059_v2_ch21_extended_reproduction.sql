-- +goose Up
CREATE TABLE extended_reproduction_schemes (
    id VARCHAR(24) NOT NULL,
    period VARCHAR(20) NOT NULL,
    accumulation_rate_i_bps BIGINT NOT NULL DEFAULT 0,
    accumulation_rate_ii_bps BIGINT NOT NULL DEFAULT 0,
    is_balanced TINYINT(1) NOT NULL DEFAULT 0,
    tick_count BIGINT NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE extended_departmental_capitals (
    id VARCHAR(24) NOT NULL,
    scheme_id VARCHAR(24) NOT NULL,
    department VARCHAR(20) NOT NULL,
    constant_pence BIGINT NOT NULL,
    variable_pence BIGINT NOT NULL,
    surplus_pence BIGINT NOT NULL,
    total_pence BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    CONSTRAINT fk_ext_dept_scheme FOREIGN KEY (scheme_id) REFERENCES extended_reproduction_schemes(id),
    UNIQUE KEY uq_ext_scheme_dept (scheme_id, department)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE reinvestments (
    id VARCHAR(24) NOT NULL,
    scheme_id VARCHAR(24) NOT NULL,
    department VARCHAR(20) NOT NULL,
    delta_constant_pence BIGINT NOT NULL,
    delta_variable_pence BIGINT NOT NULL,
    consumed_surplus_pence BIGINT NOT NULL,
    cycle_number BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    CONSTRAINT fk_reinvest_scheme FOREIGN KEY (scheme_id) REFERENCES extended_reproduction_schemes(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE multi_period_schemes (
    id VARCHAR(24) NOT NULL,
    label VARCHAR(100) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE multi_period_scheme_entries (
    id VARCHAR(24) NOT NULL,
    multi_period_id VARCHAR(24) NOT NULL,
    scheme_id VARCHAR(24) NOT NULL,
    position BIGINT NOT NULL,
    PRIMARY KEY (id),
    CONSTRAINT fk_mp_entry_mp FOREIGN KEY (multi_period_id) REFERENCES multi_period_schemes(id),
    CONSTRAINT fk_mp_entry_scheme FOREIGN KEY (scheme_id) REFERENCES extended_reproduction_schemes(id),
    UNIQUE KEY uq_mp_position (multi_period_id, position)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE composition_shifts (
    id VARCHAR(24) NOT NULL,
    multi_period_id VARCHAR(24) NOT NULL,
    cycle_number BIGINT NOT NULL,
    dept_i_composition_bps BIGINT NOT NULL,
    dept_ii_composition_bps BIGINT NOT NULL,
    social_composition_bps BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    CONSTRAINT fk_comp_shift_mp FOREIGN KEY (multi_period_id) REFERENCES multi_period_schemes(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE dept_i_growth_leads (
    id VARCHAR(24) NOT NULL,
    multi_period_id VARCHAR(24) NOT NULL,
    cycle_number BIGINT NOT NULL,
    dept_i_growth_bps BIGINT NOT NULL,
    dept_ii_growth_bps BIGINT NOT NULL,
    lead_confirmed TINYINT(1) NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    CONSTRAINT fk_growth_lead_mp FOREIGN KEY (multi_period_id) REFERENCES multi_period_schemes(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- +goose Down
DROP TABLE IF EXISTS dept_i_growth_leads;
DROP TABLE IF EXISTS composition_shifts;
DROP TABLE IF EXISTS multi_period_scheme_entries;
DROP TABLE IF EXISTS multi_period_schemes;
DROP TABLE IF EXISTS reinvestments;
DROP TABLE IF EXISTS extended_departmental_capitals;
DROP TABLE IF EXISTS extended_reproduction_schemes;
