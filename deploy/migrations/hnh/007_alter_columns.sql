ALTER TABLE hnh_data.vacancy
    DROP CONSTRAINT IF EXISTS salary_upper_bound_is_positive,
    DROP CONSTRAINT IF EXISTS salary_lower_bound_is_positive,
    ADD CONSTRAINT salary_upper_bound_is_positive CHECK (salary_upper_bound >= 0),
    ADD CONSTRAINT salary_lower_bound_is_positive CHECK (salary_lower_bound >= 0);
    

ALTER TABLE hnh_data.cv 
    DROP CONSTRAINT IF EXISTS middle_name_is_not_empty;

---- create above / drop below ----

ALTER TABLE hnh_data.vacancy
    DROP CONSTRAINT IF EXISTS salary_upper_bound_is_positive,
    DROP CONSTRAINT IF EXISTS salary_lower_bound_is_positive,
    ADD CONSTRAINT salary_upper_bound_is_positive CHECK (salary_upper_bound > 0),
    ADD CONSTRAINT salary_lower_bound_is_positive CHECK (salary_lower_bound > 0);

ALTER TABLE hnh_data.cv 
    ADD CONSTRAINT middle_name_is_not_empty CHECK (length(middle_name) > 0);
