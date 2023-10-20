# Схема базы данных

## Нормализация

### ER диаграмма

```mermaid
erDiagram
    USER {
        int id
        string email
        string password
        string first_name
        string last_name
    }


    EMPLOYER {
        int user_id
        int vacancy_id
    }

    APPLICANT {
        int user_id
        int cv_id
        string status
    }

    CV{
        int id
        int applicant_id
        string name
        string applicant_first_name
        string applicant_last_name
        date date_of_birth
        string location
        string description
        string status
        date created
        date updated
    }

    EXPERIENCE {
        int id
        int cv_id
        string organization_name
        string position
        string description
        date start
        date end
    }

    CV_LANGUAGE_ASSIGN {
        int cv_id
        int language_id
    }

    LANGUAGE{
        int id
        string name
        string level
    }

    SKILL{
        int id
        string name
    }

    CV_SKILL_ASSIGN {
        int cv_id
        int skill_id
    }

    VACANCY_SKILL_ASSIGN {
        int vacancy_id
        int skill_id
    }

    EDUCATION {
        int cv_id
        int institution_id
        int education_major_id
        string graduation_year
    }

    EDUCATION_INSTITUTION {
        int id
        string name
        string education_level
    }

    MAJOR_FIELD {
        int id
        int institution_id
        string name
    }

    INSTITUTION_FIELD_ASSIGN {
        int institution_id
        int major_field_id
    }

    VACANCY{
        int id
        int employer_id
        string name
        string description
        int salary
        string employment
        string experience
        string education_type
        string location
        date created
        date updated
    }

    RESPOND {
        int id
        int vacancy_id
        int cv_id
    }

    ORGANIZATION{
        int employer_id
        string name
        string location
        string description
    }

    EMPLOYER ||--|| ORGANIZATION : ""
    USER ||--o| APPLICANT : ""
    USER ||--o| EMPLOYER : ""
    APPLICANT ||--o{ CV : ""

    CV ||--o{ EXPERIENCE : ""

    CV }o--|| CV_LANGUAGE_ASSIGN : ""
    CV_LANGUAGE_ASSIGN }o--|| LANGUAGE : ""

    CV ||--|| EDUCATION : ""


    CV }o--|| CV_SKILL_ASSIGN : ""
    CV_SKILL_ASSIGN }o--|| SKILL : ""



    VACANCY }o--|| RESPOND : ""
    RESPOND ||--o{ CV : ""

    VACANCY }o--|| VACANCY_SKILL_ASSIGN : ""
    SKILL }o--|| VACANCY_SKILL_ASSIGN : ""

    EMPLOYER ||--o{ VACANCY : ""

    EDUCATION }o--|| EDUCATION_INSTITUTION : ""
    EDUCATION }o--|| MAJOR_FIELD : ""
    INSTITUTION_FIELD_ASSIGN ||--|{ EDUCATION_INSTITUTION : ""
    MAJOR_FIELD }|--|| INSTITUTION_FIELD_ASSIGN : ""

```

### Функциональные зависимости

#### USER

Отношение `USER`, содержит стандартную информацию о пользователе. Имеет связь 1:1 с отношениями `APPLICANT` и `EMPLOYER`, что позволяет давать различные роли одному пользователю.

```
Relation USER:
    {id} -> password, first_name, last_name
    {email} -> password, first_name, last_name
```

В отношении `USER` выполняются следующие нормальные формы:

- **1 НФ** - значения атрибутов `id`, `email`, `password`, `first_name`, `last_name` являются атомарными
- **2 НФ** - `password`, `first_name`, `last_name` функционально зависят полностью от первичного ключа `id` и потенциального ключа `email`
- **3 НФ** - среди неключевых атрибутов `email`, `password`, `first_name`, `last_name` нет функциональных зависимостей
- **НФБК** - `email` как детерминант функциональной зависимости является потенциальным ключом

---

#### APPLICANT

Отношение `APPLICANT` имеет связь 1:1 с отношением `USER` и связь 1:M с отношением `CV`. Также обладает атрибутом `status`, который говорит о статусе поиска работы соискателем.

```
Relation APPLICANT:
    {user_id} -> cv_id, status
```

В отношении `APPLICANT` выполняются следующие нормальные формы:

- **1 НФ** - значения атрибутов `cv_id`, `status` являются атомарными
- **2 НФ** - `cv_id`, `status` функционально зависят полностью от первичного ключа `user_id`
- **3 НФ** - среди неключевых атрибутов `cv_id`, `status` нет функциональных зависимостей
- **НФБК** - все детерминанты являются потенциальными ключами

---

#### EMPLOYER

Отношение `EMPLOYER` имеет связи:

- 1:1 с отношением `USER`
- 1:M с отношением `VACANCY`
- 1:1 с отношением `ORGANIZATION`.

```
Relation EMPLOYER:
    {user_id} -> vacancy_id
```

В отношении `EMPLOYER` выполняются следующие нормальные формы:

- **1 НФ** - значения атрибутов `vacancy_id` является атомарными
- **2 НФ** - `vacancy_id` функционально зависит полностью от первичного ключа `user_id`
- **3 НФ** - среди неключевых атрибутов `vacancy_id` нет функциональных зависимостей
- **НФБК** - все детерминанты являются потенциальными ключами

---

#### ORGANIZATION

Отношение `ORGANIZATION` содержит основную информацию об организации. Также содержит связь 1:1 с отношением `EMPLOYER`.

```
Relation ORGANIZATION:
    {employer_id} -> name, location, description
```

В отношении `ORGANIZATION` выполняются следующие нормальные формы:

- **1 НФ** - значения атрибутов `name`, `location`, `description` являются атомарными
- **2 НФ** - `name`, `location`, `description` функционально зависят полностью от первичного ключа `employer_id`
- **3 НФ** - среди неключевых атрибутов `name`, `location`, `description` нет функциональных зависимостей
- **НФБК** - все детерминанты являются потенциальными ключами

---

#### VACANCY

Отношение `VACANCY` содержит основную информацию о вакансии. Имеет связи:

- M:1 с отношением `EMPLOYER`
- M:N с отношением `SKILL`

```
Relation VACANCY:
    {id} -> employer_id, education_type_id, name, description, salary, employment, experience, education_type, location, created, updated
```

В отношении `VACANCY` выполняются следующие нормальные формы:

- **1 НФ** - значения атрибутов `employer_id`, `name`, `description`, `salary`, `employment`, `experience`, `education_type`, `location`, `created`, `updated` являются атомарными
- **2 НФ** - `employer_id`, `name`, `description`, `salary`, `employment`, `experience`, `education_type`, `location`, `created`, `updated` функционально зависят полностью от первичного ключа `id`
- **3 НФ** - среди неключевых атрибутов `employer_id`, `name`, `description`, `salary`, `employment`, `experience`, `education_type`, `location`, `created`, `updated` нет функциональных зависимостей
- **НФБК** - все детерминанты являются потенциальными ключами

---

#### CV

Отношение `CV` содержит основную информацию о резюме. Имеет связи:

- M:1 с отношением `APPLICANT`
- 1:M с отношением `EXPERIENCE`
- M:N с отношением `LANGUAGE`
- M:N с отношением `SKILL`
- 1:1 с отношением `EDUCATION`

```
Relation CV:
    {id} -> applicant_id, name, applicant_first_name, applicant_last_name, date_of_birth, location, description, status, created, updated
```

В отношении `CV` выполняются следующие нормальные формы:

- **1 НФ** - значения атрибутов `applicant_id`, `name`, `applicant_first_name`, `applicant_last_name`, `date_of_birth`, `location`, `description`, `status`, `created`, `updated` являются атомарными
- **2 НФ** - `applicant_id`, `name`, `applicant_first_name`, `applicant_last_name`, `date_of_birth`, `location`, `description`, `status`, `created`, `updated` функционально зависят полностью от первичного ключа `id`
- **3 НФ** - среди неключевых атрибутов `applicant_id`, `name`, `applicant_first_name`, `applicant_last_name`, `date_of_birth`, `location`, `description`, `status`, `created`, `updated` нет функциональных зависимостей
- **НФБК** - все детерминанты являются потенциальными ключами

---

#### EXPERIENCE

Отношение `EXPERIENCE` содержит основную информацию об опыте работы. Имеет связи:

- M:1 с отношением `CV`

```
Relation EXPERIENCE:
    {id} -> cv_id, organization_name, position, description, start, end
```

В отношении `EXPERIENCE` выполняются следующие нормальные формы:

- **1 НФ** - значения атрибутов `cv_id`, `organization_name`, `position`, `description`, `start`, `end` являются атомарными
- **2 НФ** - `cv_id`, `organization_name`, `position`, `description`, `start`, `end`, `created`, `updated` функционально зависят полностью от первичного ключа `id`
- **3 НФ** - среди неключевых атрибутов `cv_id`, `organization_name`, `position`, `description`, `start`, `end`, `created`, `updated` нет функциональных зависимостей
- **НФБК** - все детерминанты являются потенциальными ключами

---

#### LANGUAGE

Отношение `LANGUAGE` содержит основную информацию об уровне знания языка. Имеет связи:

- M:N с отношением `CV`

```
Relation LANGUAGE:
    {id} -> name, level
```

В отношении `LANGUAGE` выполняются следующие нормальные формы:

- **1 НФ** - значения атрибутов `name`, `level` являются атомарными
- **2 НФ** - `name`, `level` функционально зависят полностью от первичного ключа `id`
- **3 НФ** - среди неключевых атрибутов `name`, `level` нет функциональных зависимостей
- **НФБК** - все детерминанты являются потенциальными ключами

---

#### SKILL

Отношение `SKILL` содержит основную информацию о навыках. Имеет связи:

- M:N с отношением `CV`
- M:N с отношением `VACANCY`

```
Relation SKILL:
    {id} -> name
```

В отношении `SKILL` выполняются следующие нормальные формы:

- **1 НФ** - значения атрибутов `name` являются атомарными
- **2 НФ** - `name` функционально зависят полностью от первичного ключа `id`
- **3 НФ** - среди неключевых атрибутов `name` нет функциональных зависимостей
- **НФБК** - все детерминанты являются потенциальными ключами

---

#### EDUCATION

Отношение `EDUCATION` содержит основную информацию об образовании. Имеет связи:

- 1:1 с отношением `CV`
- M:1 с отношением `MAJOR_FIELD`
- M:1 с отношением `EDUCATION_INSTITUTION`

```
Relation EDUCATION:
    {cv_id} -> institution_id, education_major_id, graduation_year
```

В отношении `EDUCATION` выполняются следующие нормальные формы:

- **1 НФ** - значения атрибутов `institution_id`, `education_major_id`, `graduation_year` являются атомарными
- **2 НФ** - `institution_id`, `education_major_id`, `graduation_year` функционально зависят полностью от первичного ключа `cv_id`
- **3 НФ** - среди неключевых атрибутов `institution_id`, `education_major_id`, `graduation_year` нет функциональных зависимостей
- **НФБК** - все детерминанты являются потенциальными ключами

---

#### EDUCATION_INSTITUTION

Отношение `EDUCATION_INSTITUTION` содержит основную информацию об учебном учереждении. Имеет связи:

- 1:M с отношением `EDUCATION`
- M:N с отношением `MAJOR_FIELD`

```
Relation EDUCATION_INSTITUTION:
    {id} -> name, education_level
```

В отношении `EDUCATION_INSTITUTION` выполняются следующие нормальные формы:

- **1 НФ** - значения атрибутов `name`, `education_level` являются атомарными
- **2 НФ** - `name`, `education_level` функционально зависят полностью от первичного ключа `id`
- **3 НФ** - среди неключевых атрибутов `name`, `education_level` нет функциональных зависимостей
- **НФБК** - все детерминанты являются потенциальными ключами

---

#### MAJOR_FIELD

Отношение `MAJOR_FIELD` содержит основную информацию о специальности. Имеет связи:

- M:N с отношением `EDUCATION_INSTITUTION`
- 1:M с отношением `EDUCATION`

```
Relation MAJOR_FIELD:
    {id} -> institution_id, name
```

В отношении `MAJOR_FIELD` выполняются следующие нормальные формы:

- **1 НФ** - значения атрибутов `institution_id`, `name` являются атомарными
- **2 НФ** - `institution_id`, `name` функционально зависят полностью от первичного ключа `id`
- **3 НФ** - среди неключевых атрибутов `institution_id`, `name` нет функциональных зависимостей
- **НФБК** - все детерминанты являются потенциальными ключами

---

#### RESPOND

Отношение `RESPOND` содержит основную информацию об откликах соискателей на вакансии. Имеет связи:

- 1:M с отношением `VACANCY`
- 1:M с отношением `CV`

```
Relation RESPOND:
    {id} -> vacancy_id, cv_id
```

В отношении `RESPOND` выполняются следующие нормальные формы:

- **1 НФ** - значения атрибутов `vacancy_id`, `cv_id` являются атомарными
- **2 НФ** - `vacancy_id`, `cv_id` функционально зависят полностью от первичного ключа `id`
- **3 НФ** - среди неключевых атрибутов `vacancy_id`, `cv_id` нет функциональных зависимостей
- **НФБК** - все детерминанты являются потенциальными ключами
