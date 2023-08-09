package users

import (
	"github.com/evermos/boilerplate-go/infras"
	"github.com/rs/zerolog/log"
)

type UserRepository interface {
	GetData(filter UserFilter, page, size int) ([]UserView, error)
	CountTotalData(filter UserFilter) (int, error)
	GetProfile(uid string) (*ProfileView, error)
}

type UserRepositoryMySQL struct {
	DB *infras.MySQLConn
}

func ProvideUserRepositoryMySQL(db *infras.MySQLConn) *UserRepositoryMySQL {
	return &UserRepositoryMySQL{
		DB: db,
	}
}

func (r *UserRepositoryMySQL) GetData(filter UserFilter, page, size int) ([]UserView, error) {
	query := `
		SELECT 
			u.username,
		 	p.name,
			u.role,
			p.gender,
			p.dob,
			p.education,
			p.city,
			p.province,
			p.address,
			p.phone_number,
			s.job_role,
			s.status,
			pl.city AS placement,
			d.name AS department_name
		FROM 
			ums_users AS u
		LEFT JOIN
			ums_profiles AS p
				ON u.profile_id = p.id
		LEFT JOIN
			ums_status AS s
				ON u.status_id = s.id
		LEFT JOIN
			ums_placement AS pl
				ON u.placement_id = pl.id
		LEFT JOIN
			ums_dept AS d
				ON u.dept_id = d.id
	`

	// Add filters
	args := []interface{}{}
	if filter.Name != "" {
		if len(args) > 0 {
			query += " AND"
		} else {
			query += " WHERE"
		}
		query += " p.name LIKE ?"
		args = append(args, "%"+filter.Name+"%")
	}

	if filter.City != "" {
		if len(args) > 0 {
			query += " AND"
		} else {
			query += " WHERE"
		}
		query += " p.city LIKE ?"
		args = append(args, "%"+filter.City+"%")
	}

	if filter.Province != "" {
		if len(args) > 0 {
			query += " AND"
		} else {
			query += " WHERE"
		}
		query += " p.province LIKE ?"
		args = append(args, "%"+filter.Province+"%")
	}

	if filter.JobRole != "" {
		if len(args) > 0 {
			query += " AND"
		} else {
			query += " WHERE"
		}
		query += " s.job_role LIKE ?"
		args = append(args, "%"+filter.JobRole+"%")
	}

	if filter.Status != "" {
		if len(args) > 0 {
			query += " AND"
		} else {
			query += " WHERE"
		}
		query += " s.status LIKE ?"
		args = append(args, "%"+filter.Status+"%")
	}

	// Add pagination
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 5
	}

	query += " LIMIT ? OFFSET ?"
	offset := (page - 1) * size
	args = append(args, size, offset)

	var users []UserView
	err := r.DB.Read.Select(&users, query, args...)
	if err != nil {
		log.Error().Err(err).Msg("Failed to read users from db")
		return nil, err
	}

	return users, nil
}

func (r *UserRepositoryMySQL) CountTotalData(filter UserFilter) (int, error) {
	totalDataQuery := `
		SELECT 
			COUNT(*) 
		FROM 
			ums_users AS u
		LEFT JOIN 
			ums_profiles AS p 
				ON u.profile_id = p.id
		LEFT JOIN 
			ums_status AS s 
				ON u.status_id = s.id
		LEFT JOIN 
			ums_placement AS pl 
				ON u.placement_id = pl.id
		LEFT JOIN 
			ums_dept AS d 
				ON u.dept_id = d.id
	`

	// Add filters
	argsTotalData := []interface{}{}
	if filter.Name != "" {
		if len(argsTotalData) > 0 {
			totalDataQuery += " AND"
		} else {
			totalDataQuery += " WHERE"
		}
		totalDataQuery += " p.name LIKE ?"
		argsTotalData = append(argsTotalData, "%"+filter.Name+"%")
	}

	if filter.City != "" {
		if len(argsTotalData) > 0 {
			totalDataQuery += " AND"
		} else {
			totalDataQuery += " WHERE"
		}
		totalDataQuery += " p.city LIKE ?"
		argsTotalData = append(argsTotalData, "%"+filter.City+"%")
	}

	if filter.Province != "" {
		if len(argsTotalData) > 0 {
			totalDataQuery += " AND"
		} else {
			totalDataQuery += " WHERE"
		}
		totalDataQuery += " p.province LIKE ?"
		argsTotalData = append(argsTotalData, "%"+filter.Province+"%")
	}

	if filter.JobRole != "" {
		if len(argsTotalData) > 0 {
			totalDataQuery += " AND"
		} else {
			totalDataQuery += " WHERE"
		}
		totalDataQuery += " s.job_role LIKE ?"
		argsTotalData = append(argsTotalData, "%"+filter.JobRole+"%")
	}

	if filter.Status != "" {
		if len(argsTotalData) > 0 {
			totalDataQuery += " AND"
		} else {
			totalDataQuery += " WHERE"
		}
		totalDataQuery += " s.status LIKE ?"
		argsTotalData = append(argsTotalData, "%"+filter.Status+"%")
	}

	var totalData int
	err := r.DB.Read.Get(&totalData, totalDataQuery, argsTotalData...)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get total data")
		return 0, err
	}

	return totalData, nil
}

func (r *UserRepositoryMySQL) GetProfile(uid string) (*ProfileView, error) {
	query := `
	SELECT 
		p.name,
		u.role,
		p.gender,
		p.dob,
		p.education,
		p.city,
		p.province,
		p.address,
		p.phone_number,
		s.job_role,
		s.status,
		pl.city AS placement_city,
		d.name AS department_name
	FROM 
		ums_users AS u
	LEFT JOIN
		ums_profiles AS p
			ON u.profile_id = p.id
	LEFT JOIN
		ums_status AS s
			ON u.status_id = s.id
	LEFT JOIN
		ums_placement AS pl
			ON u.placement_id = pl.id
	LEFT JOIN
		ums_dept AS d
			ON u.dept_id = d.id
	WHERE u.id = ?
	`

	var profile ProfileView
	err := r.DB.Read.Get(&profile, query, uid)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get profile")
		return nil, err
	}
	return &profile, nil
}
