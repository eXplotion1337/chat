CREATE TABLE "public".users (
    "id" INTEGER PRIMARY KEY GENERATED BY DEFAULT AS IDENTITY,
    "nick" TEXT NOT NULL,
    "passwordHash" TEXT DEFAULT '' NOT NULL,
    "createdAt" TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL
);