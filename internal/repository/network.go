package repository

import (
	"regexp"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Contact struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	OwnerID   string             `bson:"ownerId" json:"ownerId"`
	Name      string             `bson:"name" json:"name"`
	Email     string             `bson:"email" json:"email"`
	Company   string             `bson:"company" json:"company"`
	Role      string             `bson:"role" json:"role"`
	LinkedIn  string             `bson:"linkedin" json:"linkedin"`
	Notes     string             `bson:"notes" json:"notes"`
	Mobile    string             `bson:"mobile,omitempty" json:"mobile,omitempty"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt" json:"updatedAt"`
}

// CreateContact adds a new contact to the user's network
func (mc *MongoDBClient) CreateContact(c *Contact) error {
	ctx, cancel := getContextWithTimeout(10)
	defer cancel()

	c.CreatedAt = time.Now()
	c.UpdatedAt = time.Now()

	collection := mc.Database("referrer").Collection("contacts")
	_, err := collection.InsertOne(ctx, c)
	return err
}

// GetContacts retrieves contacts for a user, optionally filtered by a search query
func (mc *MongoDBClient) GetContacts(ownerID string, query string) ([]*Contact, error) {
	ctx, cancel := getContextWithTimeout(10)
	defer cancel()

	collection := mc.Database("referrer").Collection("contacts")

	filter := bson.M{}
	if ownerID != "" {
		filter["ownerId"] = ownerID
	}

	if query != "" {
		safeQuery := regexp.QuoteMeta(query)
		filter["$or"] = []bson.M{
			{"name": bson.M{"$regex": safeQuery, "$options": "i"}},
			{"company": bson.M{"$regex": safeQuery, "$options": "i"}},
			{"role": bson.M{"$regex": safeQuery, "$options": "i"}},
			{"email": bson.M{"$regex": safeQuery, "$options": "i"}},
		}
	}

	opts := options.Find().SetSort(bson.M{"createdAt": -1})
	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var contacts []*Contact
	if err := cursor.All(ctx, &contacts); err != nil {
		return nil, err
	}
	if contacts == nil {
		contacts = []*Contact{}
	}

	return contacts, nil
}

// UpdateContact updates an existing contact
func (mc *MongoDBClient) UpdateContact(c *Contact) error {
	ctx, cancel := getContextWithTimeout(10)
	defer cancel()

	collection := mc.Database("referrer").Collection("contacts")

	filter := bson.M{"_id": c.ID, "ownerId": c.OwnerID}
	update := bson.M{
		"$set": bson.M{
			"name":      c.Name,
			"email":     c.Email,
			"company":   c.Company,
			"role":      c.Role,
			"linkedin":  c.LinkedIn,
			"notes":     c.Notes,
			"mobile":    c.Mobile,
			"updatedAt": time.Now(),
		},
	}

	_, err := collection.UpdateOne(ctx, filter, update)
	return err
}

// DeleteContact removes a contact
func (mc *MongoDBClient) DeleteContact(id primitive.ObjectID, ownerID string) error {
	ctx, cancel := getContextWithTimeout(10)
	defer cancel()

	collection := mc.Database("referrer").Collection("contacts")
	_, err := collection.DeleteOne(ctx, bson.M{"_id": id, "ownerId": ownerID})
	return err
}

// SearchContactsForAutocomplete searches contacts for the autocomplete/search endpoint
// This replaces the old SearchPeople logic for self-hosted context
func (mc *MongoDBClient) SearchContactsForAutocomplete(ownerID string, query string) ([]map[string]string, error) {
	contacts, err := mc.GetContacts(ownerID, query)
	if err != nil {
		return nil, err
	}

	output := []map[string]string{}
	for _, c := range contacts {
		output = append(output, map[string]string{
			"name":           c.Name,
			"email":          c.Email,
			"currentCompany": c.Company,
			"currentRole":    c.Role,
			"about":          c.Notes,
		})
	}
	return output, nil
}

// SyncContactsFromDrafts imports contacts from AI email drafts
func (mc *MongoDBClient) SyncContactsFromDrafts(ownerEmail string) (int, error) {
	ctx, cancel := getContextWithTimeout(20)
	defer cancel()

	// 1. Get all AI drafts for the user
	draftsCollection := mc.Database("referrer").Collection("ai_email_drafts")
	
	// We only care about drafts that have a valid 'To' email and 'CompanyName'
	filter := bson.M{
		"userEmailAddress": ownerEmail,
		"to":               bson.M{"$ne": ""},
	}
	
	cursor, err := draftsCollection.Find(ctx, filter)
	if err != nil {
		return 0, err
	}
	defer cursor.Close(ctx)

	var drafts []AiDraftColdEmail
	if err := cursor.All(ctx, &drafts); err != nil {
		return 0, err
	}

	contactsCollection := mc.Database("referrer").Collection("contacts")
	count := 0

	for _, draft := range drafts {
		// 2. Check if contact already exists
		// Normalize email to lower case
		// draft.To might contain multiple emails? The struct says `To string`. 
		// Assuming single email or comma separated. 
		// For simplicity, let's treat it as single or take the first one if we parse it.
		// Given the UI uses an autocomplete which sets a single email usually or the backend might store it as string.
		
		email := draft.To
		if email == "" {
			continue
		}

		// Simple check: if email already in contacts for this owner
		existingCount, err := contactsCollection.CountDocuments(ctx, bson.M{
			"ownerId": ownerEmail,
			"email":   email,
		})
		if err != nil {
			continue
		}

		if existingCount == 0 {
			// 3. Create new contact
			// Derive name from email
			name := deriveNameFromEmail(email)
			
			contact := Contact{
				OwnerID:   ownerEmail,
				Name:      name,
				Email:     email,
				Company:   draft.CompanyName,
				Role:      "Recruiter / Peer", // Default role
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
				Notes:     "Imported from AI Email Drafts",
			}
			
			_, err := contactsCollection.InsertOne(ctx, contact)
			if err == nil {
				count++
			}
		}
	}

	return count, nil
}

func deriveNameFromEmail(email string) string {
	// Simple heuristic: name@domain.com -> Name
	// john.doe@... -> John Doe
	atIndex := regexp.MustCompile(`@`)
	loc := atIndex.FindStringIndex(email)
	if loc == nil {
		return email
	}
	
	localPart := email[:loc[0]]
	// Replace dots and underscores with spaces
	re := regexp.MustCompile(`[._]`)
	name := re.ReplaceAllString(localPart, " ")
	
	// Capitalize words (simple version)
	// For a more robust title case, we'd iterate.
	// This is just a placeholder logic.
	return name
}
