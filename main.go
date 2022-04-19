package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"sort"
	"strings"
)

type depName string

const (
	bioDepName    depName = "Biotech"
	chemDepName   depName = "Chemistry"
	comSciDepName depName = "Engineering"
	mathDepName   depName = "Mathematics"
	physDepName   depName = "Physics"
	noneName      depName = "none"
)

var departmentNames = []depName{
	bioDepName,
	chemDepName,
	comSciDepName,
	mathDepName,
	physDepName,
}

type name struct {
	first string
	last  string
}

func (name name) getFullName() string {
	return fmt.Sprintf("%s %s", name.first, name.last)
}

func (dep depName) toString() string {
	return string(dep)
}

type application struct {
	id                 string
	name               name
	scores             scores
	departments        []depName
	assignedDepartment depName
}

func (app application) getRelevantScore() float64 {
	switch app.assignedDepartment {
	case bioDepName:
		return getBestScore(app.scores.admission, app.scores.chemistry, app.scores.physics)
	case chemDepName:
		return getBestScore(app.scores.admission, app.scores.chemistry)
	case comSciDepName:
		return getBestScore(app.scores.admission, app.scores.math, app.scores.computerScience)
	case mathDepName:
		return getBestScore(app.scores.admission, app.scores.math)
	case physDepName:
		return getBestScore(app.scores.admission, app.scores.math, app.scores.physics)
	}

	return 0.0
}

func (app application) writeToFile(file *os.File) {
	file.WriteString(fmt.Sprintf("%s %.1f\n", app.name.getFullName(), app.getRelevantScore()))
}

func getBestScore(admissionScore float64, nums ...float64) float64 {
	count := 0.0
	sum := 0.0
	for _, num := range nums {
		sum += num
		count++
	}
	return math.Max(admissionScore, sum/count)
}

func (apps applications) sort(name depName) {
	sort.Slice(apps, func(i, j int) bool {
		appI, appJ := apps[i], apps[j]
		iScores, jScores := appI.scores, appJ.scores
		switch name {
		case bioDepName:
			iScore, jScore := getBestScore(iScores.admission, iScores.physics, iScores.chemistry), getBestScore(jScores.admission, jScores.physics, jScores.chemistry)
			if iScore != jScore {
				return iScore > jScore
			}
		case chemDepName:
			iScore, jScore := getBestScore(iScores.admission, iScores.chemistry), getBestScore(jScores.admission, jScores.chemistry)
			if iScore != jScore {
				return iScore > jScore
			}
		case comSciDepName:
			iScore, jScore := getBestScore(iScores.admission, iScores.computerScience, iScores.math), getBestScore(jScores.admission, jScores.computerScience, jScores.math)
			if iScore != jScore {
				return iScore > jScore
			}
		case mathDepName:
			iScore, jScore := getBestScore(iScores.admission, iScores.math), getBestScore(jScores.admission, jScores.math)
			if iScore != jScore {
				return iScore > jScore
			}
		case physDepName:
			iScore, jScore := getBestScore(iScores.admission, iScores.physics, iScores.math), getBestScore(jScores.admission, jScores.physics, jScores.math)
			if iScore != jScore {
				return iScore > jScore
			}
		}

		return apps[i].name.getFullName() < apps[j].name.getFullName()
	})
}

func (apps applications) writeToFiles() {
	for _, name := range departmentNames {
		apps.sort(name)
		fileName := fmt.Sprintf("%s.txt", strings.ToLower(name.toString()))
		file, err := os.Create(fileName)
		if err != nil {
			panic(err)
		}
		defer file.Close()
		for _, app := range apps {
			if app.assignedDepartment == name {
				app.writeToFile(file)
			}
		}
	}
}

type scores struct {
	physics         float64
	chemistry       float64
	math            float64
	computerScience float64
	admission       float64
}

type applications []application

var allApplications applications

func main() {
	var maxApplicants int
	fmt.Scan(&maxApplicants)
	readFile()

	assignDepartments(maxApplicants)

	allApplications.writeToFiles()
}

func readFile() {
	file, err := os.Open("./applicants.txt")
	if err != nil {
		log.Fatal(err)
	}

	var apps applications
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var firstName, lastName string
		var physicsScore, chemistryScore, mathScore, computerScienceScore, admissionScore float64
		var department1, department2, department3 depName

		_, err := fmt.Sscanf(scanner.Text(), "%s %s %f %f %f %f %f %s %s %s", &firstName, &lastName, &physicsScore, &chemistryScore, &mathScore, &computerScienceScore, &admissionScore, &department1, &department2, &department3)
		if err != nil {
			panic(err)
		}

		apps = append(apps, application{
			name: name{first: firstName, last: lastName},
			scores: scores{
				physics:         physicsScore,
				chemistry:       chemistryScore,
				math:            mathScore,
				computerScience: computerScienceScore,
				admission:       admissionScore,
			},
			departments: []depName{
				department1,
				department2,
				department3,
			},
			assignedDepartment: noneName,
		})
	}

	allApplications = apps
}

func countAssigned(dep depName) int {
	count := 0
	for _, app := range allApplications {
		if app.assignedDepartment == dep {
			count = count + 1
		}
	}

	return count
}

func assignDepartments(max int) {
	for flow := 0; flow < 3; flow++ {
		assignDepartmentsByFlow(max, flow)
	}
}

func assignDepartmentsByFlow(max int, flow int) {
	for _, dep := range departmentNames {
		allApplications.sort(dep)
		for index, app := range allApplications {
			if app.departments[flow] != dep {
				continue
			}
			if app.departments[flow] == dep && app.assignedDepartment == noneName && countAssigned(dep) < max {
				allApplications[index].assignedDepartment = dep
			}
		}
	}
}
