const form = document.getElementById("recommendation-form");
const resultDiv = document.getElementById("result");
const machineSlots = document.querySelectorAll(".machine-slot");
const exerciseOptions = document.getElementById("exercise-options");
const availableCount = document.getElementById("available-count");
const occupiedCount = document.getElementById("occupied-count");
const availableList = document.getElementById("available-list");
const logoutButton = document.getElementById("logout-button");
const currentUserSpan = document.getElementById("current-user");
const machineMessage = document.getElementById("machine-message");
// --- NEW VARIABLES FOR ROUTINES ---
const routineSelect = document.getElementById("routine-select");
const startRoutineBtn = document.getElementById("start-routine-btn");
const routineProgressDiv = document.getElementById("routine-progress");
const currentRoutineExerciseSpan = document.getElementById("current-routine-exercise");
const nextExerciseBtn = document.getElementById("next-exercise-btn");
const exerciseInput = document.getElementById("exercise");

let allRoutines = [];      // This is where we'll store the complete routines returned by the database
let currentRoutine = [];   // List of exercise names in the active routine
let currentExerciseIndex = 0;

form.addEventListener("submit", async function (event) {
    event.preventDefault();

    const exercise = document.getElementById("exercise").value.trim();
    resultDiv.classList.remove("result-success", "result-error", "result-neutral");

    if (!exercise) {
        resultDiv.classList.add("result-error");
        resultDiv.innerHTML = `
            <h2>Recommendation</h2>
            <p>Please enter an exercise name.</p>
        `;
        return;
    }

    try {
        const response = await fetch(`/recommendation?exercise=${encodeURIComponent(exercise)}`);
        const data = await response.json();

        if (!response.ok) {
            resultDiv.classList.add("result-error");
            resultDiv.innerHTML = `
                <h2>Recommendation</h2>
                <p>${data.message || "Could not get recommendation."}</p>
            `;
            return;
        }

        resultDiv.classList.add("result-success");
        resultDiv.innerHTML = `
            <h2>Recommendation</h2>
            <p><strong>Requested Exercise:</strong> ${data.requested_exercise}</p>
            <p><strong>Recommended Exercise:</strong> ${data.recommended_exercise}</p>
            <p><strong>Muscle Group:</strong> ${data.muscle_group}</p>
            <p><strong>Machine:</strong> ${data.machine}</p>
        `;
    } catch (error) {
        resultDiv.classList.add("result-error");
        resultDiv.innerHTML = `
            <h2>Recommendation</h2>
            <p>An unexpected error occurred while fetching the recommendation.</p>
        `;
    }
});

async function loadMachines() {
    try {
        const response = await fetch("/machines");
        const machines = await response.json();

        const availableMachines = machines.filter(machine => machine.is_available);
        const occupiedMachines = machines.filter(machine => !machine.is_available);

        availableCount.textContent = availableMachines.length;
        occupiedCount.textContent = occupiedMachines.length;

        if (availableMachines.length > 0) {
            availableList.textContent = availableMachines.map(machine => machine.name).join(", ");
        } else {
            availableList.textContent = "No machines available";
        }

        machineSlots.forEach(slot => {
            slot.classList.remove("available", "occupied");

            const machineName = slot.dataset.machineName;
            const machine = machines.find(m => m.name === machineName);

            const statusText = slot.querySelector(".machine-status");
            if (statusText) statusText.textContent = "Unknown";

            if (machine) {
                slot.dataset.machineId = machine.id;
                slot.dataset.available = machine.is_available;

                if (machine.is_available) {
                    slot.classList.add("available");

                    if (statusText) {
                        statusText.textContent = "Available";
                    }
                } else {
                    slot.classList.add("occupied");

                    if (statusText) {

                        if (machine.occupied_until) {

                            const occupiedUntil = new Date(machine.occupied_until);
                            const now = new Date();

                            const diffMs = occupiedUntil - now;

                            if (diffMs > 0) {

                                const totalSeconds = Math.floor(diffMs / 1000);

                                const minutes = Math.floor(totalSeconds / 60);
                                const seconds = totalSeconds % 60;

                                statusText.textContent =
                                    `Occupied (${minutes}m ${seconds}s)`;

                            } else {

                                statusText.textContent = "Occupied";
                            }

                        } else {

                            statusText.textContent = "Occupied";
                        }
                    }
                }
            }
        });
    } catch (error) {
        console.error("Could not load machines:", error);
    }
}

async function toggleMachineAvailability(machineId, currentAvailability) {
    try {
        const newAvailability = !currentAvailability;

        const response = await fetch("/machines/update-availability-post", {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify({
                id: Number(machineId),
                available: newAvailability
            })
        });

        if (!response.ok) {
            const errorMessage = await response.text();

            if (machineMessage) {
                machineMessage.textContent = errorMessage || "Could not update machine availability.";
            }

            return;
        }

        await loadMachines();
        if (machineMessage) {
            machineMessage.textContent = "";
        }
    } catch (error) {
        console.error("Error updating machine availability:", error);
    }
}

machineSlots.forEach(slot => {
    slot.addEventListener("click", async function () {
        const machineId = this.dataset.machineId;
        const currentAvailability = this.dataset.available === "true";

        if (!machineId) {
            return;
        }

        await toggleMachineAvailability(machineId, currentAvailability);
    });
});

async function loadExercises() {
    try {
        const response = await fetch("/exercises");
        const exercises = await response.json();

        exerciseOptions.innerHTML = "";

        exercises.forEach(exercise => {
            const option = document.createElement("option");
            option.value = exercise.name;
            exerciseOptions.appendChild(option);
        });
    } catch (error) {
        console.error("Could not load exercises:", error);
    }
}


if (logoutButton) {
    logoutButton.addEventListener("click", async function () {
        try {
            const response = await fetch("/api/logout", {
                method: "POST"
            });

            if (!response.ok) {
                console.error("Could not logout");
                return;
            }

            window.location.href = "/login";
        } catch (error) {
            console.error("Unexpected error during logout:", error);
        }
    });
}

async function loadCurrentUser() {
    if (!currentUserSpan) {
        return;
    }

    try {
        const response = await fetch("/api/me");

        if (!response.ok) {
            window.location.href = "/login";
            return;
        }

        const user = await response.json();

        currentUserSpan.textContent = `Logged in as ${user.name} · ${user.role}`;
        
        //We pass the actual ID returned by the session API
        loadRoutines(user.id);
    } catch (error) {
        console.error("Could not load current user:", error);
        window.location.href = "/login";
    }
}
// --- LOGIC FOR ROUTINES ---
// Function that now receives the actual ID of the logged-in user
async function loadRoutines(userId) {
    try {
        const response = await fetch(`/routines?userId=${userId}`);
        if (!response.ok) return;
        
        allRoutines = await response.json();
        
        // We clear the selector just in case and leave the default option
        routineSelect.innerHTML = '<option value="">Select a Routine</option>';
        
        allRoutines.forEach(routine => {
            const option = document.createElement("option");
            option.value = routine.id;
            option.textContent = routine.name;
            routineSelect.appendChild(option);
        });
    } catch (error) {
        console.error("Could not load routines:", error);
    }
}

// Enable the “Start” button only when a routine is selected
routineSelect.addEventListener("change", (e) => {
    startRoutineBtn.disabled = e.target.value === "";
});

// When you click “Start Routine” (We map the actual exercises from the JSON in the database)
startRoutineBtn.addEventListener("click", () => {
    const selectedRoutineId = parseInt(routineSelect.value);
    
    // We look for the selected routine in our local list 'allRoutines'
    const selectedRoutine = allRoutines.find(r => r.id === selectedRoutineId);
    
    if (!selectedRoutine || !selectedRoutine.exercises || selectedRoutine.exercises.length === 0) {
        alert("This routine has no exercises assigned yet.");
        return;
    }

    // We extract only the names of the exercises, preserving the order from the database
    currentRoutine = selectedRoutine.exercises.map(re => re.exercise.name);
    currentExerciseIndex = 0;
    
    routineProgressDiv.style.display = "block";
    updateRoutineUI();
});

// Button to advance through the training circuit
nextExerciseBtn.addEventListener("click", () => {
    currentExerciseIndex++;
    if (currentExerciseIndex < currentRoutine.length) {
        updateRoutineUI();
    } else {
        routineProgressDiv.style.display = "none";
        alert("Routine Finished! Great Job!");
        currentRoutine = [];
    }
});

// Automatically populates the recommendation search bar with the current exercise
function updateRoutineUI() {
    const nextExerciseName = currentRoutine[currentExerciseIndex];
    currentRoutineExerciseSpan.textContent = nextExerciseName;
    
    // We enter the exercise into your recommendation form and click submit
    exerciseInput.value = nextExerciseName;
    form.dispatchEvent(new Event('submit'));
}

// Add the initial call alongside the others
loadRoutines();

loadMachines();
loadExercises();
loadCurrentUser();

setInterval(loadMachines, 1000);