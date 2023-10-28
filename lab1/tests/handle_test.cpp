#include <handle.h>
#include <gtest/gtest.h>

TEST(HandleTests, Test_One)
{
    int in = 1;
    int expected = 1;
    int actual = handle(in);
    ASSERT_EQ(expected, actual);
}

int main(int argc, char** argv)
{
    ::testing::InitGoogleTest(&argc, argv);
    return RUN_ALL_TESTS();
}
