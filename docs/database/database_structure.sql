CREATE TABLE "categories" (
  "category_id" serial PRIMARY KEY NOT NULL,
  "name" text UNIQUE NOT NULL,
  "min_age" integer NOT NULL
);

CREATE TABLE "complaints" (
  "complaint_id" serial PRIMARY KEY NOT NULL,
  "complainer_id" integer NOT NULL,
  "reported_id" integer NOT NULL,
  "reason" text NOT NULL,
  "description" text NOT NULL,
  "created_at" timestamp DEFAULT (now())
);

CREATE TABLE "lessons" (
  "lesson_id" serial PRIMARY KEY NOT NULL,
  "student_id" integer NOT NULL,
  "teacher_id" integer NOT NULL,
  "category_id" integer NOT NULL,
  "schedule_time_id" integer UNIQUE NOT NULL,
  "price" integer NOT NULL DEFAULT 0,
  "status_id" integer,
  "state_machine_item_id" integer NOT NULL
);

CREATE TABLE "reviews" (
  "review_id" serial PRIMARY KEY NOT NULL,
  "teacher_id" integer NOT NULL,
  "student_id" integer NOT NULL,
  "category_id" integer NOT NULL,
  "skill_id" integer NOT NULL,
  "rate" smallint NOT NULL,
  "comment" text NOT NULL DEFAULT (''::text)
);

CREATE TABLE "schedule_times" (
  "schedule_time_id" serial PRIMARY KEY NOT NULL,
  "teacher_id" integer NOT NULL,
  "datetime" timestamp NOT NULL,
  "is_available" boolean NOT NULL DEFAULT true
);

CREATE TABLE "skills" (
  "skill_id" serial PRIMARY KEY NOT NULL,
  "teacher_id" integer NOT NULL,
  "category_id" integer NOT NULL,
  "video_card_link" text,
  "about" text,
  "rate" doubleprecision NOT NULL DEFAULT 0,
  "total_rate_score" integer NOT NULL DEFAULT 0,
  "reviews_count" integer NOT NULL DEFAULT 0,
  "is_active" boolean NOT NULL DEFAULT false
);

CREATE TABLE "state_machines" (
  "state_machine_id" serial PRIMARY KEY NOT NULL,
  "name" text UNIQUE NOT NULL,
  "start_state_id" integer NOT NULL
);

CREATE TABLE "state_machines_items" (
  "item_id" serial PRIMARY KEY NOT NULL,
  "state_machine_id" integer NOT NULL,
  "state_id" integer NOT NULL,
  "created_at" timestamp DEFAULT (now())
);

CREATE TABLE "state_transitions" (
  "transition_id" serial PRIMARY KEY NOT NULL,
  "state_machine_id" integer NOT NULL,
  "current_state_id" integer NOT NULL,
  "next_state_id" integer NOT NULL
);

CREATE TABLE "states" (
  "state_id" serial PRIMARY KEY NOT NULL,
  "name" text UNIQUE NOT NULL
);

CREATE TABLE "statuses" (
  "status_id" serial PRIMARY KEY NOT NULL,
  "name" text UNIQUE NOT NULL
);

CREATE TABLE "teachers" (
  "teacher_id" serial PRIMARY KEY NOT NULL,
  "user_id" integer UNIQUE NOT NULL,
  "rate" doubleprecision NOT NULL DEFAULT 0,
  "total_rate_score" integer NOT NULL DEFAULT 0,
  "reviews_count" integer NOT NULL DEFAULT 0
);

CREATE TABLE "users" (
  "user_id" serial PRIMARY KEY NOT NULL,
  "kratos_identity_id" UUID UNIQUE NOT NULL,
  "email" text UNIQUE NOT NULL,
  "name" text NOT NULL,
  "surname" text NOT NULL,
  "password" text NOT NULL,
  "registration_date" timestamp DEFAULT (now()),
  "birthdate" date NOT NULL,
  "avatar" text NOT NULL DEFAULT (''::text),
  "is_admin" boolean NOT NULL DEFAULT false
);

CREATE UNIQUE INDEX "unique_review" ON "reviews" ("teacher_id", "student_id", "category_id", "skill_id");

CREATE UNIQUE INDEX "unique_teacher_schedule_time" ON "schedule_times" ("teacher_id", "datetime");

CREATE UNIQUE INDEX "unique_teacher_category" ON "skills" ("teacher_id", "category_id");

CREATE UNIQUE INDEX "state_transitions_state_machine_id_current_state_id_next_st_key" ON "state_transitions" ("state_machine_id", "current_state_id", "next_state_id");

ALTER TABLE "complaints" ADD CONSTRAINT "complaints_complainer_id_fkey" FOREIGN KEY ("complainer_id") REFERENCES "users" ("user_id") ON DELETE CASCADE;

ALTER TABLE "complaints" ADD CONSTRAINT "complaints_reported_id_fkey" FOREIGN KEY ("reported_id") REFERENCES "users" ("user_id") ON DELETE CASCADE;

ALTER TABLE "lessons" ADD CONSTRAINT "lessons_category_id_fkey" FOREIGN KEY ("category_id") REFERENCES "categories" ("category_id");

ALTER TABLE "lessons" ADD CONSTRAINT "lessons_schedule_time_id_fkey" FOREIGN KEY ("schedule_time_id") REFERENCES "schedule_times" ("schedule_time_id");

ALTER TABLE "lessons" ADD CONSTRAINT "lessons_state_machine_item_id_fkey" FOREIGN KEY ("state_machine_item_id") REFERENCES "state_machines_items" ("item_id");

ALTER TABLE "lessons" ADD CONSTRAINT "lessons_status_id_fkey" FOREIGN KEY ("status_id") REFERENCES "statuses" ("status_id");

ALTER TABLE "lessons" ADD CONSTRAINT "lessons_student_id_fkey" FOREIGN KEY ("student_id") REFERENCES "users" ("user_id");

ALTER TABLE "lessons" ADD CONSTRAINT "lessons_teacher_id_fkey" FOREIGN KEY ("teacher_id") REFERENCES "teachers" ("teacher_id");

ALTER TABLE "reviews" ADD CONSTRAINT "reviews_category_id_fkey" FOREIGN KEY ("category_id") REFERENCES "categories" ("category_id") ON DELETE CASCADE;

ALTER TABLE "reviews" ADD CONSTRAINT "reviews_skill_id_fkey" FOREIGN KEY ("skill_id") REFERENCES "skills" ("skill_id") ON DELETE CASCADE;

ALTER TABLE "reviews" ADD CONSTRAINT "reviews_student_id_fkey" FOREIGN KEY ("student_id") REFERENCES "users" ("user_id") ON DELETE CASCADE;

ALTER TABLE "reviews" ADD CONSTRAINT "reviews_teacher_id_fkey" FOREIGN KEY ("teacher_id") REFERENCES "teachers" ("teacher_id") ON DELETE CASCADE;

ALTER TABLE "schedule_times" ADD CONSTRAINT "schedule_times_teacher_id_fkey" FOREIGN KEY ("teacher_id") REFERENCES "teachers" ("teacher_id") ON DELETE CASCADE;

ALTER TABLE "skills" ADD CONSTRAINT "skills_category_id_fkey" FOREIGN KEY ("category_id") REFERENCES "categories" ("category_id") ON DELETE CASCADE;

ALTER TABLE "skills" ADD CONSTRAINT "skills_teacher_id_fkey" FOREIGN KEY ("teacher_id") REFERENCES "teachers" ("teacher_id") ON DELETE CASCADE;

ALTER TABLE "state_machines_items" ADD CONSTRAINT "state_machines_items_state_id_fkey" FOREIGN KEY ("state_id") REFERENCES "states" ("state_id");

ALTER TABLE "state_machines_items" ADD CONSTRAINT "state_machines_items_state_machine_id_fkey" FOREIGN KEY ("state_machine_id") REFERENCES "state_machines" ("state_machine_id");

ALTER TABLE "state_machines" ADD CONSTRAINT "state_machines_start_state_id_fkey" FOREIGN KEY ("start_state_id") REFERENCES "states" ("state_id");

ALTER TABLE "state_transitions" ADD CONSTRAINT "state_transitions_current_state_id_fkey" FOREIGN KEY ("current_state_id") REFERENCES "states" ("state_id");

ALTER TABLE "state_transitions" ADD CONSTRAINT "state_transitions_next_state_id_fkey" FOREIGN KEY ("next_state_id") REFERENCES "states" ("state_id");

ALTER TABLE "state_transitions" ADD CONSTRAINT "state_transitions_state_machine_id_fkey" FOREIGN KEY ("state_machine_id") REFERENCES "state_machines" ("state_machine_id");

ALTER TABLE "teachers" ADD CONSTRAINT "teachers_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "users" ("user_id") ON DELETE CASCADE;
