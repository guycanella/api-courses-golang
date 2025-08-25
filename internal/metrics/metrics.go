package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var CoursesCreatedTotal = promauto.NewCounter(prometheus.CounterOpts{
	Name: "total_created_courses",
	Help: "Total number of courses created successfully",
})
