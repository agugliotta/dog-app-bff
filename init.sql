DROP TABLE IF EXISTS pets CASCADE;
DROP TABLE IF EXISTS breeds CASCADE;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE breeds (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    temperament TEXT,
    origin VARCHAR(255),
    slug TEXT UNIQUE NOT NULL

);

CREATE TABLE pets (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    birth DATE NOT NULL,
    breed_id UUID NOT NULL REFERENCES breeds(id)
);

-- Inserciones de razas
INSERT INTO breeds (name, temperament, origin, slug) VALUES
('Golden Retriever', 'Friendly, Intelligent, Devoted', 'Scotland', 'golden-retriever'),
('Labrador Retriever', 'Outgoing, Even-tempered, Gentle', 'Canada', 'labrador-retriever'),
('German Shepherd', 'Intelligent, Obedient, Courageous', 'Germany', 'german-shepherd'),
('Bulldog', 'Docile, Willful, Friendly', 'England', 'bulldog'),
('Poodle', 'Intelligent, Proud, Active', 'Germany/France', 'poodle'),
('Beagle', 'Gentle, Excitable, Determined', 'England', 'beagle'),
('Rottweiler', 'Good-natured, Devoted, Confident', 'Germany', 'rottweiler'),
('Dachshund', 'Clever, Stubborn, Playful', 'Germany', 'dachshund'),
('Siberian Husky', 'Outgoing, Mischievous, Alert', 'Siberia', 'siberian-husky'),
('Boxer', 'Playful, Loyal, Energetic', 'Germany', 'boxer'),
('Shih Tzu', 'Playful, Outgoing, Affectionate', 'Tibet', 'shih-tzu'),
('Great Dane', 'Gentle, Friendly, Dependable', 'Germany', 'great-dane'),
('Chihuahua', 'Graceful, Charming, Big-Personality', 'Mexico', 'chihuahua'),
('Pomeranian', 'Playful, Lively, Extroverted', 'Germany/Poland', 'pomeranian'),
('Australian Shepherd', 'Intelligent, Active, Protective', 'United States', 'australian-shepherd'),
('Yorkshire Terrier', 'Spunky, Bold, Energetic', 'England', 'yorkshire-terrier'),
('Cavalier King Charles Spaniel', 'Affectionate, Gentle, Graceful', 'United Kingdom', 'cavalier-king-charles-spaniel'),
('Doberman Pinscher', 'Fearless, Alert, Loyal', 'Germany', 'doberman-pinscher'),
('Shetland Sheepdog', 'Playful, Intelligent, Responsive', 'Scotland', 'shetland-sheepdog'),
('Mastiff', 'Courageous, Dignified, Good-natured', 'England', 'mastiff'),
('Bichon Frise', 'Playful, Charming, Gentle', 'Spain/Belgium', 'bichon-frise'),
('Newfoundland', 'Sweet-tempered, Devoted, Patient', 'Canada', 'newfoundland'),
('Basenji', 'Intelligent, Independent, Quiet', 'Congo', 'basenji'),
('Greyhound', 'Gentle, Independent, Noble', 'Middle East/Egypt', 'greyhound'),
('Affenpinscher', 'Loyal, Curious, Stubborn', 'Germany', 'affenpinscher'),
('French Bulldog', 'Affectionate, Calm, Playful', 'France', 'french-bulldog'),
('Bernese Mountain Dog', 'Good-natured, Calm, Affectionate', 'Switzerland', 'bernese-mountain-dog'),
('Alaskan Malamute', 'Affectionate, Loyal, Playful', 'United States', 'alaskan-malamute'),
('Pug', 'Charming, Mischievous, Stubborn', 'China', 'pug'),
('Dalmatian', 'Outgoing, Playful, Energetic', 'Croatia', 'dalmatian');

-- Inserciones de ejemplo de mascotas
INSERT INTO pets (name, birth, breed_id) VALUES
('Buddy', '2022-05-10', (SELECT id FROM breeds WHERE name = 'Golden Retriever')),
('Max', '2023-01-20', (SELECT id FROM breeds WHERE name = 'Bulldog')),
('Lucy', '2024-03-15', (SELECT id FROM breeds WHERE name = 'Siberian Husky'));
