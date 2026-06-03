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

const routineSelect = document.getElementById("routine-select");
const startRoutineBtn = document.getElementById("start-routine-btn");
const routineProgressDiv = document.getElementById("routine-progress");
const currentRoutineExerciseSpan = document.getElementById("current-routine-exercise");
const nextExerciseBtn = document.getElementById("next-exercise-btn");
const exerciseInput = document.getElementById("exercise");
const routinePreview = document.getElementById("routine-preview");
const routinePreviewList = document.getElementById("routine-preview-list");

let allRoutines = [];
let globalExercises = [];
let currentRoutine = [];
let currentExerciseIndex = 0;
let userOccupiedMachineId = null;

/* =========================
   TOAST NOTIFICATIONS
========================= */
function showToast(message, type = 'success') {
    const container = document.getElementById('toast-container');
    if (!container) {
        console.error("No se encontró el toast-container en el HTML");
        return;
    }

    const toast = document.createElement('div');
    toast.className = `toast ${type}`;
    toast.textContent = message;

    container.appendChild(toast);

    setTimeout(() => {
        toast.classList.add('fade-out');
        toast.addEventListener('animationend', () => {
            toast.remove();
        });
    }, 3500);
}

function showConfirmModal(message, title = "Confirm Action") {
    return new Promise((resolve) => {
        const overlay = document.getElementById('custom-confirm-overlay');
        if (!overlay) {
            console.error("No se encontró el overlay del modal en HTML");
            return resolve(false);
        }

        const titleEl = overlay.querySelector('.confirm-title');
        const messageEl = document.getElementById('confirm-message');
        const cancelBtn = document.getElementById('confirm-cancel-btn');
        const okBtn = document.getElementById('confirm-ok-btn');

        titleEl.textContent = title;
        messageEl.textContent = message;

        overlay.classList.remove('hidden');

        const cleanup = () => {
            overlay.classList.add('hidden');
            cancelBtn.removeEventListener('click', onCancel);
            okBtn.removeEventListener('click', onOk);
        };

        const onCancel = () => { cleanup(); resolve(false); };
        const onOk = () => { cleanup(); resolve(true); };

        cancelBtn.addEventListener('click', onCancel);
        okBtn.addEventListener('click', onOk);
    });
}

function populateExerciseDropdown(dropdownId, exercisesList) {
    const dropdown = document.getElementById(dropdownId);
    if (!dropdown) return;

    dropdown.innerHTML = '<option value="" selected disabled>Choose an exercise...</option>';

    exercisesList.forEach(ex => {
        const option = document.createElement("option");
        option.value = ex.id;
        option.textContent = ex.name;
        dropdown.appendChild(option);
    });
}

/* =========================
   RECOMMENDATION FORM
========================= */

form.addEventListener("submit", async function (event) {

    event.preventDefault();

    const exercise = exerciseInput.value.trim();

    resultDiv.classList.remove(
        "result-success",
        "result-error",
        "result-neutral"
    );

    if (!exercise) {

        resultDiv.classList.add("result-error");

        resultDiv.innerHTML = `
            <h2>Recommendation</h2>
            <p>Please enter an exercise name.</p>
        `;

        return;
    }

    try {

        const response = await fetch(
            `/recommendation?exercise=${encodeURIComponent(exercise)}`
        );

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

            <p>
                <strong>Requested Exercise:</strong>
                ${data.requested_exercise}
            </p>

            <p>
                <strong>Recommended Exercise:</strong>
                ${data.recommended_exercise}
            </p>

            <p>
                <strong>Muscle Group:</strong>
                ${data.muscle_group}
            </p>

            <p>
                <strong>Machine:</strong>
                ${data.machine}
            </p>
        `;

    } catch (error) {

        resultDiv.classList.add("result-error");

        resultDiv.innerHTML = `
            <h2>Recommendation</h2>
            <p>
                An unexpected error occurred while fetching
                the recommendation.
            </p>
        `;
    }
});

/* =========================
   LOAD MACHINES
========================= */

async function loadMachines() {

    try {

        const response = await fetch("/machines");

        const machines = await response.json();

        const availableMachines =
            machines.filter(machine => machine.is_available);

        const occupiedMachines =
            machines.filter(machine => !machine.is_available);

        availableCount.textContent = availableMachines.length;

        occupiedCount.textContent = occupiedMachines.length;

        if (availableMachines.length > 0) {

            availableList.textContent =
                availableMachines
                    .map(machine => machine.name)
                    .join(", ");

        } else {

            availableList.textContent =
                "No machines available";
        }

        machineSlots.forEach(slot => {

            slot.classList.remove(
                "available",
                "occupied"
            );

            const machineName =
                slot.dataset.machineName;

            const machine =
                machines.find(m => m.name === machineName);

            const statusText =
                slot.querySelector(".machine-status");

            if (statusText) {
                statusText.textContent = "Unknown";
            }

            if (machine) {

                slot.dataset.machineId = machine.id;

                slot.dataset.available =
                    machine.is_available;

                if (machine.is_available) {

                    slot.classList.add("available");

                    if (statusText) {

                        statusText.textContent =
                            "Available";
                    }

                } else {

                    slot.classList.add("occupied");

                    if (statusText) {

                        if (machine.occupied_until) {

                            const occupiedUntil =
                                new Date(machine.occupied_until);

                            const now = new Date();

                            const diffMs =
                                occupiedUntil - now;

                            if (diffMs > 0) {

                                const totalSeconds =
                                    Math.floor(diffMs / 1000);

                                const minutes =
                                    Math.floor(totalSeconds / 60);

                                const seconds =
                                    totalSeconds % 60;

                                statusText.textContent =
                                    `Occupied (${minutes}m ${seconds}s)`;

                            } else {

                                statusText.textContent =
                                    "Occupied";
                            }

                        } else {

                            statusText.textContent =
                                "Occupied";
                        }
                    }
                }
            }
        });

    } catch (error) {

        console.error(
            "Could not load machines:",
            error
        );
    }
}

/* =========================
   MACHINE AVAILABILITY
========================= */

async function toggleMachineAvailability(
    machineId,
    currentAvailability
) {

    try {

        const newAvailability =
            !currentAvailability;

        machineMessage.classList.remove(
            "success",
            "error",
            "warning",
            "info"
        );

        machineMessage.textContent =
            "Updating machine status...";

        machineMessage.classList.add("info");

        const response = await fetch(
            "/machines/update-availability-post",
            {
                method: "POST",

                headers: {
                    "Content-Type": "application/json"
                },

                body: JSON.stringify({
                    id: Number(machineId),
                    available: newAvailability
                })
            }
        );

        const responseText =
            await response.text();

        if (!response.ok) {

            machineMessage.classList.remove("info");

            machineMessage.classList.add("error");

            machineMessage.textContent =
                responseText ||
                "Could not update machine availability.";

            setTimeout(() => {

                machineMessage.textContent = "";

                machineMessage.classList.remove("error");

            }, 5000);

            return;
        }

        await loadMachines();

        machineMessage.classList.remove("info");

        machineMessage.classList.add("success");

        if (newAvailability) {

            userOccupiedMachineId = Number(machineId);

            machineMessage.textContent =
                "Machine successfully occupied.";

        } else {

            userOccupiedMachineId = null;

            machineMessage.textContent =
                "Machine successfully released.";
        }

        setTimeout(() => {

            machineMessage.textContent = "";

            machineMessage.classList.remove(
                "success",
                "error",
                "warning",
                "info"
            );

        }, 3500);

    } catch (error) {

        console.error(
            "Error updating machine availability:",
            error
        );

        machineMessage.classList.remove("info");

        machineMessage.classList.add("error");

        machineMessage.textContent =
            "Unexpected error while updating machine.";

        setTimeout(() => {

            machineMessage.textContent = "";

            machineMessage.classList.remove("error");

        }, 5000);
    }
}

/* =========================
   MACHINE CLICK EVENTS
========================= */

machineSlots.forEach(slot => {

    slot.addEventListener("click", async function () {

        const machineId =
            this.dataset.machineId;

        const currentAvailability =
            this.dataset.available === "true";

        if (!machineId) {
            return;
        }


        if (
            currentAvailability &&
            userOccupiedMachineId !== null &&
            userOccupiedMachineId !== Number(machineId)
        ) {

            machineMessage.classList.remove(
                "success",
                "info"
            );

            machineMessage.classList.add("warning");

            machineMessage.textContent =
                "You already occupy another machine. Release it before selecting a new one.";

            setTimeout(() => {

                machineMessage.textContent = "";

                machineMessage.classList.remove("warning");

            }, 5000);

            return;
        }

        await toggleMachineAvailability(
            machineId,
            currentAvailability
        );
    });
});

/* =========================
   LOAD EXERCISES
========================= */

async function loadExercises() {
    try {
        const response = await fetch("/exercises");

        globalExercises = await response.json();

        if (exerciseOptions) {
            exerciseOptions.innerHTML = "";
            globalExercises.forEach(ex => {
                const datalistOption = document.createElement("option");
                datalistOption.value = ex.name;
                exerciseOptions.appendChild(datalistOption);
            });
        }

        populateExerciseDropdown("new-routine-exercise-dropdown", globalExercises);
        populateExerciseDropdown("edit-routine-exercise-dropdown", globalExercises);

    } catch (error) {
        console.error("Could not load exercises:", error);
    }
}

// ==========================================
// --- LIVE SEARCH ---
// ==========================================
const searchCreateInput = document.getElementById('search-create-exercise');
if (searchCreateInput) {
    searchCreateInput.addEventListener('input', (e) => {
        const searchTerm = e.target.value.toLowerCase();
        const filteredExercises = globalExercises.filter(ex =>
            ex.name.toLowerCase().includes(searchTerm)
        );
        populateExerciseDropdown("new-routine-exercise-dropdown", filteredExercises);
    });
}

const searchEditInput = document.getElementById('search-edit-exercise');
if (searchEditInput) {
    searchEditInput.addEventListener('input', (e) => {
        const searchTerm = e.target.value.toLowerCase();
        const filteredExercises = globalExercises.filter(ex =>
            ex.name.toLowerCase().includes(searchTerm)
        );
        populateExerciseDropdown("edit-routine-exercise-dropdown", filteredExercises);
    });
}

/* =========================
   LOGOUT
========================= */

if (logoutButton) {

    logoutButton.addEventListener(
        "click",
        async function () {

            try {

                const response = await fetch(
                    "/api/logout",
                    {
                        method: "POST"
                    }
                );

                if (!response.ok) {

                    console.error("Could not logout");

                    return;
                }

                window.location.href = "/login";

            } catch (error) {

                console.error(
                    "Unexpected error during logout:",
                    error
                );
            }
        }
    );
}

/* =========================
   CURRENT USER
========================= */

async function loadCurrentUser() {

    if (!currentUserSpan) {
        return;
    }

    try {

        const response =
            await fetch("/api/me");

        if (!response.ok) {

            window.location.href = "/login";

            return;
        }

        const user =
            await response.json();

        currentUserSpan.textContent =
            `Logged in as ${user.name} · ${user.role}`;


        window.loggedInUserId = user.id;

        loadRoutines(user.id);

    } catch (error) {

        console.error(
            "Could not load current user:",
            error
        );

        window.location.href = "/login";
    }
}

/* =========================
   ROUTINES
========================= */

async function loadRoutines(userId) {
    try {
        const response = await fetch(`/routines?userId=${userId}`);
        if (!response.ok) return;

        const data = await response.json();
        allRoutines = data || [];

        routineSelect.innerHTML = '<option value="">Select a Routine</option>';

        const deleteRoutineSelect = document.getElementById("delete-routine-select");
        if (deleteRoutineSelect) {
            deleteRoutineSelect.innerHTML = '<option value="">Select a routine to delete</option>';
        }

        const editRoutineSelect = document.getElementById("edit-routine-select");
        if (editRoutineSelect) {
            editRoutineSelect.innerHTML = '<option value="">Select a routine...</option>';
        }

        allRoutines.forEach(routine => {
            const option = document.createElement("option");
            option.value = routine.id;
            option.textContent = routine.name;
            routineSelect.appendChild(option);

            if (deleteRoutineSelect) {
                const delOption = document.createElement("option");
                delOption.value = routine.id;
                delOption.textContent = routine.name;
                deleteRoutineSelect.appendChild(delOption);
            }

            if (editRoutineSelect) {
                const editOption = document.createElement("option");
                editOption.value = routine.id;
                editOption.textContent = routine.name;
                editRoutineSelect.appendChild(editOption);
            }
        });
    } catch (error) {
        console.error("Could not load routines:", error);
    }
}

/* =========================
   ROUTINE SELECT
========================= */

routineSelect.addEventListener("change", (e) => {
    const hasValue = e.target.value !== "";
    startRoutineBtn.disabled = !hasValue;

    const deleteBtn = document.getElementById("delete-routine-btn");
    if (deleteBtn) deleteBtn.disabled = !hasValue;

    if (hasValue && routinePreview && routinePreviewList) {
        const selectedRoutineId = parseInt(e.target.value);
        const selectedRoutine = allRoutines.find(r => r.id === selectedRoutineId);

        routinePreviewList.innerHTML = "";

        if (selectedRoutine && selectedRoutine.exercises && selectedRoutine.exercises.length > 0) {
            selectedRoutine.exercises.forEach(ex => {
                const li = document.createElement("li");
                li.textContent = ex.exercise.name;
                routinePreviewList.appendChild(li);
            });
        } else {
            routinePreviewList.innerHTML = "<li><em>No exercises assigned yet.</em></li>";
        }

        routinePreview.style.display = "block";
    } else if (routinePreview) {
        routinePreview.style.display = "none";
    }
});

/* =========================
   START ROUTINE
========================= */

startRoutineBtn.addEventListener("click", () => {

    const selectedRoutineId =
        parseInt(routineSelect.value);

    const selectedRoutine =
        allRoutines.find(
            r => r.id === selectedRoutineId
        );

    if (
        !selectedRoutine ||
        !selectedRoutine.exercises ||
        selectedRoutine.exercises.length === 0
    ) {

        showToast("This routine has no exercises assigned yet.", "warning");

        return;
    }

    currentRoutine =
        selectedRoutine.exercises.map(
            re => re.exercise.name
        );

    currentExerciseIndex = 0;

    routinePreview.style.display = "none";

    routineProgressDiv.style.display = "block";

    updateRoutineUI();
});

/* =========================
   NEXT EXERCISE
========================= */

nextExerciseBtn.addEventListener("click", () => {

    currentExerciseIndex++;

    if (
        currentExerciseIndex <
        currentRoutine.length
    ) {

        updateRoutineUI();

    } else {

        routineProgressDiv.style.display = "none";

        showToast("Routine Finished! Great Job!", "success");

        currentRoutine = [];
    }
});

/* =========================
   UPDATE ROUTINE UI
========================= */

function updateRoutineUI() {

    const nextExerciseName =
        currentRoutine[currentExerciseIndex];

    currentRoutineExerciseSpan.textContent =
        nextExerciseName;

    exerciseInput.value =
        nextExerciseName;

    form.dispatchEvent(new Event("submit"));
}

/* =========================
   INITIAL LOAD
========================= */


loadMachines();

loadExercises();

loadCurrentUser();

setInterval(loadMachines, 1000);

// ==========================================
// --- SIDE MENU NAVIGATION ---
// ==========================================
const navDashboard = document.getElementById('nav-dashboard');
const navRoutines = document.getElementById('nav-routines');
const navMachines = document.getElementById('nav-machines');


// sections of the Dashboard
const dashHeader = document.getElementById('dash-header');
const dashTopGrid = document.getElementById('dash-top-grid');
const dashGymSection = document.getElementById('dash-gym-section');

// Routines view
const routinesView = document.getElementById('routines-view');
const machinesView = document.getElementById('machines-view');

if (navRoutines && navDashboard && navMachines) {

    navRoutines.addEventListener('click', (e) => {

        e.preventDefault();

        navDashboard.classList.remove('active');
        navMachines.classList.remove('active');

        navRoutines.classList.add('active');

        dashHeader.style.display = 'none';
        dashTopGrid.style.display = 'none';
        dashGymSection.style.display = 'none';

        machinesView.style.display = 'none';

        routinesView.style.display = 'block';
    });

    navMachines.addEventListener('click', (e) => {

        e.preventDefault();

        navDashboard.classList.remove('active');
        navRoutines.classList.remove('active');

        navMachines.classList.add('active');

        dashHeader.style.display = 'none';
        dashTopGrid.style.display = 'none';
        dashGymSection.style.display = 'none';

        routinesView.style.display = 'none';

        machinesView.style.display = 'block';
    });

    navDashboard.addEventListener('click', (e) => {

        e.preventDefault();

        navRoutines.classList.remove('active');
        navMachines.classList.remove('active');

        navDashboard.classList.add('active');

        routinesView.style.display = 'none';
        machinesView.style.display = 'none';

        dashHeader.style.display = '';
        dashTopGrid.style.display = '';
        dashGymSection.style.display = '';
    });
}

// ==========================================
// --- LOGIC OF THE TOUCH SELECTOR (CREATE) ---
// ==========================================
const btnAddExercise = document.getElementById('btn-add-exercise');
const exerciseDropdown = document.getElementById('new-routine-exercise-dropdown');
const selectedExercisesList = document.getElementById('selected-exercises-list');

if (btnAddExercise && exerciseDropdown && selectedExercisesList) {
    btnAddExercise.addEventListener('click', () => {
        const selectedOption = exerciseDropdown.options[exerciseDropdown.selectedIndex];
        if (!selectedOption.value || selectedOption.disabled) return;

        const exerciseId = selectedOption.value;
        const exerciseName = selectedOption.textContent;

        const li = document.createElement('li');
        li.className = 'exercise-list-item';
        li.dataset.id = exerciseId;

        li.innerHTML = `
            <span>${exerciseName}</span>
            <div class="exercise-actions">
                <button type="button" class="exercise-action-btn btn-up" title="Move Up">⬆️</button>
                <button type="button" class="exercise-action-btn btn-down" title="Move Down">⬇️</button>
                <button type="button" class="exercise-action-btn btn-remove" title="Remove">❌</button>
            </div>
        `;

        li.querySelector('.btn-remove').addEventListener('click', () => li.remove());

        li.querySelector('.btn-up').addEventListener('click', () => {
            const prev = li.previousElementSibling;
            if (prev) li.parentNode.insertBefore(li, prev);
        });

        li.querySelector('.btn-down').addEventListener('click', () => {
            const next = li.nextElementSibling;
            if (next) li.parentNode.insertBefore(next, li);
        });

        selectedExercisesList.appendChild(li);
        exerciseDropdown.selectedIndex = 0;
    });
}
// ==========================================
// --- SUBMIT THE NEW ROUTINE FORM ---
// ==========================================
const createRoutineForm = document.getElementById('create-routine-form');

if (createRoutineForm) {
    createRoutineForm.addEventListener('submit', async (e) => {
        e.preventDefault();

        const nameInput = document.getElementById('new-routine-name').value.trim();

        // --- (VALIDATIONS) ---

        if (!nameInput) {
            showToast("Routine name cannot be empty.", "warning");
            return;
        }

        const nameExists = allRoutines.some(r => r.name.toLowerCase() === nameInput.toLowerCase());
        if (nameExists) {
            showToast("You already have a routine with this name. Please choose a different one.", "warning");
            return;
        }

        const listItems = document.querySelectorAll('#selected-exercises-list .exercise-list-item');
        const selectedExerciseIds = Array.from(listItems).map(li => parseInt(li.dataset.id));

        if (selectedExerciseIds.length === 0) {
            showToast("Please select at least one exercise to create a routine.", "warning");
            return;
        }

        // --- END OF VALIDATIONS ---

        const userId = window.loggedInUserId || 1;

        try {
            const response = await fetch('/routines/create', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({
                    user_id: userId,
                    name: nameInput,
                    exercise_ids: selectedExerciseIds
                })
            });

            if (response.ok) {
                showToast("Routine created successfully!", "success");
                createRoutineForm.reset();
                document.getElementById('selected-exercises-list').innerHTML = '';
                loadRoutines(userId);

                const navDashboard = document.getElementById('nav-dashboard');
                if (navDashboard) navDashboard.click();
            } else {
                showToast("Error creating routine.", "error");
            }
        } catch (error) {
            console.error("Error:", error);
        }
    });
}

// ==========================================
// --- DELETE SELECTED ROUTINE ---
// ==========================================
const deleteRoutineSelect = document.getElementById("delete-routine-select");
const deleteRoutineBtn = document.getElementById("delete-routine-btn");

if (deleteRoutineSelect && deleteRoutineBtn) {
    deleteRoutineSelect.addEventListener("change", (e) => {
        deleteRoutineBtn.disabled = e.target.value === "";
    });

    deleteRoutineBtn.addEventListener("click", async () => {
        const selectedRoutineId = parseInt(deleteRoutineSelect.value);
        if (!selectedRoutineId) return;

        const userConfirmed = await showConfirmModal("Are you sure you want to delete this routine?", "Delete Routine");
        if (!userConfirmed) {
            return;
        }

        const userId = window.loggedInUserId || 1;

        try {
            const response = await fetch(`/routines/delete?id=${selectedRoutineId}&userId=${userId}`, {
                method: "DELETE"
            });

            if (response.ok) {
                showToast("Routine deleted successfully!", "success");

                loadRoutines(userId);

                deleteRoutineBtn.disabled = true;
            } else {
                showToast("Error deleting routine. You might not be authorized.", "error");
            }
        } catch (error) {
            console.error("Error deleting routine:", error);
        }
    });
}

// ==========================================
// --- EDIT AN EXISTING ROUTINE ---
// ==========================================
const editRoutineSelect = document.getElementById("edit-routine-select");
const editRoutineForm = document.getElementById("edit-routine-form");
const editRoutineNameInput = document.getElementById("edit-routine-name");

const btnAddEditExercise = document.getElementById('btn-add-edit-exercise');
const editExerciseDropdown = document.getElementById('edit-routine-exercise-dropdown');
const editSelectedExercisesList = document.getElementById('edit-selected-exercises-list');

function addExerciseToEditList(exerciseId, exerciseName) {
    const li = document.createElement('li');
    li.className = 'exercise-list-item';
    li.dataset.id = exerciseId;

    li.innerHTML = `
        <span>${exerciseName}</span>
        <div class="exercise-actions">
            <button type="button" class="exercise-action-btn btn-up" title="Move Up">⬆️</button>
            <button type="button" class="exercise-action-btn btn-down" title="Move Down">⬇️</button>
            <button type="button" class="exercise-action-btn btn-remove" title="Remove">❌</button>
        </div>
    `;

    li.querySelector('.btn-remove').addEventListener('click', () => li.remove());

    li.querySelector('.btn-up').addEventListener('click', () => {
        const prev = li.previousElementSibling;
        if (prev) li.parentNode.insertBefore(li, prev);
    });

    li.querySelector('.btn-down').addEventListener('click', () => {
        const next = li.nextElementSibling;
        if (next) li.parentNode.insertBefore(next, li);
    });

    editSelectedExercisesList.appendChild(li);
}

if (btnAddEditExercise && editExerciseDropdown && editSelectedExercisesList) {
    btnAddEditExercise.addEventListener('click', () => {
        const selectedOption = editExerciseDropdown.options[editExerciseDropdown.selectedIndex];
        if (!selectedOption.value || selectedOption.disabled) return;

        addExerciseToEditList(selectedOption.value, selectedOption.textContent);
        editExerciseDropdown.selectedIndex = 0;
    });
}

if (editRoutineSelect && editRoutineForm) {
    editRoutineSelect.addEventListener("change", (e) => {
        const selectedRoutineId = parseInt(e.target.value);

        if (!selectedRoutineId) {
            editRoutineForm.style.display = "none";
            return;
        }

        const routineToEdit = allRoutines.find(r => r.id === selectedRoutineId);
        if (!routineToEdit) return;

        editRoutineNameInput.value = routineToEdit.name;

        editSelectedExercisesList.innerHTML = '';

        if (routineToEdit.exercises) {
            routineToEdit.exercises.forEach(ex => {
                const optionMatch = Array.from(editExerciseDropdown.options).find(opt => parseInt(opt.value) === ex.exercise_id);
                const exName = optionMatch ? optionMatch.textContent : `Exercise ${ex.exercise_id}`;

                addExerciseToEditList(ex.exercise_id, exName);
            });
        }

        editRoutineForm.style.display = "flex";
    });

    editRoutineForm.addEventListener("submit", async (e) => {
        e.preventDefault();

        const routineId = parseInt(editRoutineSelect.value);
        const nameInput = editRoutineNameInput.value.trim();

        const listItems = document.querySelectorAll('#edit-selected-exercises-list .exercise-list-item');
        const selectedExerciseIds = Array.from(listItems).map(li => parseInt(li.dataset.id));

        if (!nameInput) {
            showToast("Routine name cannot be empty.", "warning");
            return;
        }

        const nameExists = allRoutines.some(r => r.name.toLowerCase() === nameInput.toLowerCase() && r.id !== routineId);
        if (nameExists) {
            showToast("Another routine already uses this name. Please choose a different one.", "warning");
            return;
        }

        if (selectedExerciseIds.length === 0) {
            showToast("Please select at least one exercise.", "warning");
            return;
        }

        const userId = window.loggedInUserId || 1;

        try {
            const response = await fetch("/routines/update", {
                method: "PUT",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({
                    routine_id: routineId,
                    user_id: userId,
                    name: nameInput,
                    exercise_ids: selectedExerciseIds
                })
            });

            if (response.ok) {
                showToast("Routine updated successfully!", "success");
                editRoutineForm.reset();
                editSelectedExercisesList.innerHTML = '';
                editRoutineForm.style.display = "none";
                editRoutineSelect.value = "";

                loadRoutines(userId);

                const navDashboard = document.getElementById('nav-dashboard');
                if (navDashboard) navDashboard.click();
            } else {
                showToast("Error updating routine.", "error");
            }
        } catch (error) {
            console.error("Error updating routine:", error);
        }
    });
}